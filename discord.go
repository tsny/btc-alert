package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tsny/btc-alert/coinbase"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	"github.com/tsny/btc-alert/eps"
	"github.com/tsny/btc-alert/yahoo"
)

// CryptoBot is a service that communicates with discord and holds onto alerts
// that are created for discord users
type CryptoBot struct {
	ds        *discordgo.Session
	alerts    map[string]priceAlert
	channelId string
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

// todo: pass in the channel id, don't ref config all the time?
// todo: return err
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
	cb := &CryptoBot{ds: dg, channelId: conf.Discord.ChannelID}
	dg.AddHandler(cb.OnNewMessage)
	dg.AddHandler(cb.OnDisconnect)
	return cb
}

// SubscribeUserToPriceTarget alerts a user when a security hits a specific price target
// relative to the price the security was at when the user first subscribed
func (cb *CryptoBot) SubscribeUserToPriceTarget(userID string, target float64, p *eps.Publisher) {
	startedBelow := p.GetPrice() < target
	x := priceAlert{userID, p, target, p.CurrentCandle.Current, true, startedBelow}
	str := "Subbing %s to %s price point %.4f | Current: %.4f\n"
	fmt.Printf(str, userID, p.Ticker, target, p.GetPrice())
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

// SubscribeToTicker adds a ticker to the general watchlist
func (cb *CryptoBot) SubscribeToTicker(ticker string, p *eps.Publisher) {
	_ = newListener(p, conf.Intervals, conf.Thresholds)
	// PublisherMap[ticker] = p
}

// GetTopGainers outputs a table of the top gainers in the market today
func (cb *CryptoBot) GetTopGainers(gainers bool) {
	str := &strings.Builder{}
	data := yahoo.GetTopMoversAsArray(gainers)
	// Have to truncate, too many chars for a message
	data = data[0:10]
	table := tablewriter.NewWriter(str)
	table.SetHeader(yahoo.GetTableHeader())
	table.AppendBulk(data)
	table.SetCenterSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.Render()
	out := "```" + str.String() + "```"
	println("size of table " + strconv.Itoa(len(out)))
	if len(out) > 2000 {
		return
	}
	cb.SendMessage(out, "", false)
}

// SendMessage sends a discord message with an optional mention
// TODO: Could change these into options for UserID and TTS
func (cb *CryptoBot) SendMessage(str string, userID string, tts bool) (*discordgo.Message, error) {
	if userID == "" {
		return cb.ds.ChannelMessageSend(conf.Discord.ChannelID, str)
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
	return cb.ds.ChannelMessageSendComplex(conf.Discord.ChannelID, &msg)
}

// OnDisconnect logs whenever we disconnect (for debugging)
func (cb *CryptoBot) OnDisconnect(s *discordgo.Session, hb *discordgo.Disconnect) {
	logrus.Infof("Bot disconnected")
}

func (cb *CryptoBot) SendGraph(content string, reader io.Reader) {
	file := &discordgo.File{Name: "test.png", Reader: reader}
	msg := &discordgo.MessageSend{Content: content, Files: []*discordgo.File{file}}
	cb.ds.ChannelMessageSendComplex(cb.channelId, msg)
}

// OnNewMessage function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (cb *CryptoBot) OnNewMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	logrus.Infof("Processing msg '%s' from '%s'", m.Content, m.Author.Username)

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
		cb.GetTopGainers(true)
		return
	}
	if parts[0] == "losers" {
		cb.GetTopGainers(false)
		return
	}

	if len(parts) < 2 {
		return
	}

	if parts[0] == "whois" {
		t := parts[1]
		println("Getting summary for " + t)
		sum := yahoo.GetSummary(t)
		if sum == "" {
			println("Couldn't get summary for ticker " + t)
			return
		}
		// Discord don't like mo than 2000 chars
		if len(sum) > 2000 {
			sum = sum[0:1999]
		}
		cb.SendMessage(sum, "", false)
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
			cb.SendMessage("Following "+ticker, "", false)
			return
		}

		price, err := strconv.ParseFloat(parts[2], 3)
		if err != nil {
			return
		}
		cb.SubscribeUserToPriceTarget(m.Author.ID, price, pub)
		str := "Subbed %s to %s price point %.4f -- Current: %.4f"
		discordMessage := fmt.Sprintf(str, m.Author.Username, pub.Ticker, price, pub.GetPrice())
		cb.SendMessage(discordMessage, "", false)
	}

	if parts[0] == "get" {
		// Small delay, todo: get rid
		if pub.CurrentCandle == nil {
			time.Sleep(2 * time.Second)
		}
		cb.SendMessage(pub.CurrentCandle.String(), "", false)
		return
	}

	if parts[0] == "trade" {
		cdl := pub.CurrentCandle
		fee := cdl.Current * .01
		str := "%s -- $%.2f -- Fee: $%.2f -- 2%% Gain: $%.2f ($%.2f)"
		str = fmt.Sprintf(str, cdl.Ticker, cdl.Current, fee, fee*2, fee*2+cdl.Current)
		cb.SendMessage(str, "", false)
	}

	if parts[0] == "stat" {
		d := coinbase.Get24Hour(ticker)
		str := "24 Hour Status: %s -- High: $%s | Low: $%s | Open $%"
		str = fmt.Sprintf(str, ticker, d.High, d.Low, d.Open)
		cb.SendMessage(str, "", false)
	}
}
