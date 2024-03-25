package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"btc-alert/eps"
	"btc-alert/yahoo"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

// TODO: We should have crypto bot subscribe to events rather than the
// events in files like listener.go directly call crypto bot, that way
// if discord is inactive we don't have errors

// CryptoBot is a service that communicates with discord and holds onto alerts
// that are created for discord users
type CryptoBot struct {
	session         *discordgo.Session
	serverChannelID string
	alertEveryone   bool
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
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Infof("error creating Discord session: %v", err)
		return nil
	}
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsDirectMessages + discordgo.IntentGuildMessages)
	err = session.Open()
	if err != nil {
		log.Infof("error opening connection: %v", err)
		return nil
	}

	log.Infof("Connected to Discord server")
	cb := &CryptoBot{session: session, serverChannelID: conf.Discord.ChannelID}
	session.AddHandler(cb.OnNewMessage)
	cb.alertEveryone = conf.Discord.AlertEveryone

	for _, v := range conf.PercentageChanges {
		println(v.PercentChange)
	}

	go func() {
		userID := conf.Discord.UsersToNotify[0]
		channel, err := session.UserChannelCreate(userID)
		if err != nil {
			panic(err)
		}

		for {
			pub, ok := findPublisher("btc")
			if !ok {
				return
			}
			candle := *pub.Candle
			time.Sleep(time.Hour * 6)
			_, err = session.ChannelMessageSend(channel.ID, pub.Candle.Diff(&candle))
			if err != nil {
				panic(err)
			}
		}
	}()

	return cb
}

// OnNewMessage function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (cb *CryptoBot) OnNewMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	log.Infof("%v %v %v [%v]: %v", m.Message.ID, m.Message.Type, m.Author.ID, m.Author.Username, m.Content)

	msg := strings.TrimSpace(m.Content)
	if len(msg) == 0 {
		return
	}
	parts := strings.Split(msg, " ")

	if len(parts) == 0 {
		return
	}
	ticker := ""
	if len(parts) > 1 {
		ticker = strings.ToLower(parts[1])
	}

	switch parts[0] {

	case "whois", "get":
		cb.handleGet(ticker, m)
		return

	case "track":
		if len(parts) <= 2 {
			_, _ = cb.SendUserMessage(m.Author, "Invalid; usage: track btc 500")
			return
		}
		pub, ok := findPublisher(ticker)
		if !ok {
			_, _ = cb.SendUserMessage(m.Author, "ticker %v not found", ticker)
			return
		}
		// todo; need to track all listeners in case we have dupes
		cl := NewChangeListener(pub)
		chgAmount, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			chgAmount = 500
		}
		suffix := ""
		if chgAmount < 1 {
			cl.RegisterPercentListener(m.Author.ID, chgAmount)
			suffix = "%"
		} else {
			cl.Register(m.Author.ID, chgAmount)
		}
		_, _ = cb.SendUserMessage(m.Author, "Will alert you when `%v` price (`%v`) changes by `%v%v`",
			ticker, pub.Candle.Price, chgAmount, suffix)

	case "chart", "graph":
		if ticker == "" {
			_, _ = cb.SendUserMessage(m.Author, "ticker %v not found", ticker)
			return
		}
	default:
		_, _ = cb.SendUserMessage(m.Author, "What?")
	}
}

// SubscribeUserToPriceTarget alerts a user when a security hits a specific price target
// relative to the price the security was at when the user first subscribed
func (cb *CryptoBot) SubscribeUserToPriceTarget(userID string, target float64, p *eps.Publisher) {
	startedBelow := p.Price() < target
	x := priceAlert{userID, p, target, p.Candle.Price, true, startedBelow}
	str := "Subbing %s to %s price point %.4f | Current: %.4f\n"
	log.Infof(str, userID, p.Ticker, target, p.Price())
	f := func(p *eps.Publisher, candle *eps.Candlestick, completed bool) {
		if !x.active {
			return
		}
		str := fmt.Sprintf("%s Price Target %.4f Reached", p.Ticker, target)
		if startedBelow && candle.Price > target {
			cb.SendMessage(str, userID)
			x.active = false
		} else if !startedBelow && candle.Price < target {
			cb.SendMessage(str, userID)
			x.active = false
		}
	}
	p.RegisterPriceUpdateListener(f)
}

// SubscribeToTicker adds a ticker to the general watchlist
func (cb *CryptoBot) SubscribeToTicker(ticker string, p *eps.Publisher) {
	// _ = newListener(p, conf.Intervals, conf.Thresholds)
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
	cb.SendMessage(out, "")
}

// Sends a generalized message, used for alerts, 'ats' everyone if enabled
func (cb *CryptoBot) SendGeneralMessage(str string) (*discordgo.Message, error) {
	return cb.SendMessage(str, "")
}

func (cb *CryptoBot) SendAlertableMessage(str string) (*discordgo.Message, error) {
	user := ""
	if cb.alertEveryone {
		user = "everyone"
	}
	return cb.SendMessage(str, user)
}

// SendMessage sends a discord message with an optional mention
// TODO: Could change these into options for UserID and TTS
func (cb *CryptoBot) SendMessage(str string, userID string) (*discordgo.Message, error) {
	channel, err := cb.session.UserChannelCreate(userID)
	if err != nil {
		log.Errorf("err creating user channel for %v: %v", userID, err)
		return nil, err
	}
	return cb.session.ChannelMessageSend(channel.ID, str)
}

func (cb *CryptoBot) SendUserMessage(user *discordgo.User, str string, args ...interface{}) (*discordgo.Message, error) {
	return cb.SendMessage(fmt.Sprintf(str, args...), user.ID)
}

func (cb *CryptoBot) SendGraph(content string, reader io.Reader) {
	file := &discordgo.File{Name: "test.png", Reader: reader}
	msg := &discordgo.MessageSend{Content: content, Files: []*discordgo.File{file}}
	cb.session.ChannelMessageSendComplex(cb.serverChannelID, msg)
}

func (cb *CryptoBot) handleGet(ticker string, m *discordgo.MessageCreate) {
	pub, ok := findPublisher(ticker)
	if !ok {
		return
	}
	log.Infof("%v", pub.String())
	_, _ = cb.session.ChannelMessageSend(m.ChannelID, pub.String())
}
