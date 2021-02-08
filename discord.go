package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tsny/btc-alert/coinbase"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	"github.com/tsny/btc-alert/eps"
	"github.com/tsny/btc-alert/yahoo"
)

type CryptoBot struct {
	ds     *discordgo.Session
	alerts map[string]priceAlert
}

type priceAlert struct {
	requester     string
	publisher     *eps.Publisher
	targetPrice   float64
	startingPrice float64
	active        bool
	startedBelow  bool
}

var cryptoBot *CryptoBot

// SubscribeUserToPriceTarget alerts a user when a security hits a specific price target
// relative to the price the security was at when the user first subscribed
func (cb *CryptoBot) SubscribeUserToPriceTarget(userID string, target float64, p *eps.Publisher) {
	startedBelow := p.CurrentCandle.Current < target
	x := priceAlert{userID, p, target, p.CurrentCandle.Current, true, startedBelow}
	fmt.Printf("Subscribing %s to %s price point %.4f\n", userID, p.Ticker, target)
	f := func(p *eps.Publisher, candle eps.Candlestick) {
		if !x.active {
			return
		}
		str := fmt.Sprintf("%s Price Target %.4f Reached", p.Ticker, target)
		if startedBelow && candle.Current > target {
			cb.SendMessage(str, userID, true)
			x.active = false
		} else if !startedBelow && candle.Current < target {
			cb.SendMessage(str, userID, true)
			x.active = false
		}
	}
	p.Subscribe(f)
}

func (cb *CryptoBot) SubscribeToTicker(ticker string, p *eps.Publisher) {
	_ = newListener(p, conf.Intervals, conf.Thresholds)
	// PublisherMap[ticker] = p
}

// GetTopGainers outputs a table of the top gainers in the market today
func (cb *CryptoBot) GetTopGainers() {
	str := &strings.Builder{}
	data := yahoo.GetTopGainersAsArray()
	// Have to truncate, too many chars for a message
	data = data[0:9]
	table := tablewriter.NewWriter(str)
	table.SetHeader(yahoo.GetTableHeader())
	table.AppendBulk(data)
	table.SetCenterSeparator("|")
	table.Render()
	out := "```" + str.String() + "```"
	cb.SendMessage(out, "", false)
}

// SendMessage sends a discord message with an optional mention
func (cb *CryptoBot) SendMessage(str string, userID string, tts bool) {

	if userID == "" {
		cb.ds.ChannelMessageSend(conf.Discord.ChannelID, str)
		return
	}

	mention := userID
	var users []string
	if userID == "everyone" {
		mention = " @everyone"
	} else {
		mention = " <@" + userID + ">"
		users = append(users, userID)
	}
	msg := discordgo.MessageSend{
		Content: str + mention,
		TTS:     tts,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Users: users,
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeEveryone},
		},
	}
	_, err := cb.ds.ChannelMessageSendComplex(conf.Discord.ChannelID, &msg)
	if err != nil {
		println(err.Error())
	}
}

func initBot(token string) *CryptoBot {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil
	}
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return nil
	}
	println("Connected to Discord server")
	b := &CryptoBot{ds: dg}
	dg.AddHandler(b.OnNewMessage)
	return b
}

// OnNewMessage function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (cb *CryptoBot) OnNewMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	msg := strings.TrimSpace(m.Content)
	if msg[0] != '!' {
		return
	}
	parts := strings.Split(msg[1:], " ")

	// Testing funcs
	if parts[0] == "atme" {
		cb.SendMessage("test", m.Author.ID, false)
		return
	}

	if parts[0] == "atall" {
		cb.SendMessage("test", "everyone", false)
		return
	}

	if parts[0] == "gainers" {
		cb.GetTopGainers()
	}

	if len(parts) < 2 {
		return
	}

	ticker := strings.ToUpper(parts[1])
	cryptoTicker := ticker
	if i := strings.Index(ticker, "-"); i == -1 {
		cryptoTicker = ticker + "-USD"
	}
	pub, ok := PublisherMap[ticker]
	if ok {
		println("Found existing yahoo publisher for " + ticker)
	} else {
		pub, ok = PublisherMap[cryptoTicker]
		if ok {
			println("Found existing crypto publisher for " + cryptoTicker)
		} else {
			if yahoo.GetPrice(ticker) > 0 {
				pub = eps.New(yahoo.GetPrice, ticker, "Yahoo", true, 5)
				PublisherMap[ticker] = pub
				println("Made new publisher for subscriber -- " + ticker)
			} else {
				println("Could not find valid ticker for " + ticker)
				return
			}
		}
	}

	if parts[0] == "sub" {
		if len(parts) < 3 {
			cb.SubscribeToTicker(ticker, pub)
			return
		}

		price, err := strconv.ParseFloat(parts[2], 3)
		if err != nil {
			return
		}
		cb.SubscribeUserToPriceTarget(m.Author.ID, price, pub)
		str := "Subscribing %s to %s price point %.4f"
		discordMessage := fmt.Sprintf(str, m.Author.Username, pub.Ticker, price)
		cb.SendMessage(discordMessage, "", false)
	}

	if parts[0] == "get" {
		cb.SendMessage(pub.CurrentCandle.String(), "", false)
	}

	if parts[0] == "trade" {
		cdl := pub.CurrentCandle
		fee := cdl.Current * .01
		str := fmt.Sprintf("%s -- $%.2f -- Fee: $%.2f -- 2%% Gain: $%.2f", cdl.Ticker, cdl.Current, fee, fee*2)
		cb.SendMessage(str, "", false)
	}

	if parts[0] == "stat" {
		d := coinbase.Get24Hour(ticker)
		str := fmt.Sprintf("24 Hour Status: %s -- High: $%s | Low: $%s | Open $%s", ticker, d.High, d.Low, d.Open)
		cb.SendMessage(str, "", false)
	}
}
