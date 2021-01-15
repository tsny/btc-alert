package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tsny/btc-alert/coinbase"
	"github.com/tsny/btc-alert/eps"
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

// SubscribeUser is
func (cb *CryptoBot) SubscribeUser(userID string, target float64, p *eps.Publisher) {
	startedBelow := p.CurrentCandle.Current < target
	x := priceAlert{userID, p, target, p.CurrentCandle.Current, true, startedBelow}
	fmt.Printf("Subscribing %s to %s price point %.2f\n", userID, p.Source, target)
	f := func(p *eps.Publisher, candle eps.Candlestick) {
		if !x.active {
			return
		}
		str := fmt.Sprintf("%s Price Target %.2f Reached", p.Source, target)
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

// SendMessage sends a discord message with an optional mention
func (cb *CryptoBot) SendMessage(str string, userMention string, tts bool) {
	if userMention != "" {
		if userMention == "everyone" {
			str += " @" + userMention
		} else {
			msg := discordgo.MessageSend{
				Content: str + " @" + userMention,
				TTS:     tts,
				AllowedMentions: &discordgo.MessageAllowedMentions{
					Users: []string{"84090395092353024"},
					Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeEveryone},
				},
			}
			cb.ds.ChannelMessageSendComplex(conf.Discord.ChannelID, &msg)
		}
	}
	cb.ds.ChannelMessageSend(conf.Discord.ChannelID, str)
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

	if len(parts) < 2 {
		return
	}

	ticker := strings.ToUpper(parts[1])
	if i := strings.Index(ticker, "-"); i == -1 {
		ticker = ticker + "-USD"
	}
	convTicker := coinbase.Source(ticker)
	pub, ok := PublisherMap[convTicker]
	if !ok {
		println("Couldn't find publisher for " + ticker)
		return
	}

	if parts[0] == "sub" {
		if len(parts) < 3 {
			return
		}
		price, err := strconv.ParseFloat(parts[2], 3)
		if err != nil {
			return
		}
		cb.SubscribeUser(m.Author.Username, price, pub)
		str := "Subscribing %s to %s price point %.2f"
		discordMessage := fmt.Sprintf(str, m.Author.Username, pub.Source, price)
		cb.SendMessage(discordMessage, "", false)
	}

	if parts[0] == "get" {
		cb.SendMessage(pub.CurrentCandle.String(), "", false)
	}

	if parts[0] == "trade" {
		cdl := pub.CurrentCandle
		fee := cdl.Current * .01
		str := fmt.Sprintf("%s -- $%.2f -- Fee: $%.2f -- 2%% Gain: $%.2f", cdl.Source, cdl.Current, fee, fee*2)
		cb.SendMessage(str, "", false)
	}

	if parts[0] == "stat" {
		d := convTicker.Get24Hour()
		str := fmt.Sprintf("24 Hour Status: %s -- High: $%s | Low: $%s | Open $%s", ticker, d.High, d.Low, d.Open)
		cb.SendMessage(str, "", false)
	}
}
