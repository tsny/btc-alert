package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"btc-alert/priceTracking"

	log "github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"

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
		log.Infof("error creating Discord session,", err)
		return nil
	}
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		log.Infof("error opening connection,", err)
		return nil
	}

	log.Infof("Connected to Discord server")
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
	log.Infof(str, userID, p.Ticker, target, p.GetPrice())
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
	p.RegisterSubscriber(f)
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
	log.Infof("Bot disconnected")
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
	log.Infof("Processing msg '%s' from '%s'", m.Content, m.Author.Username)

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
	sec, pub := lookupService.FindSecurityByNameOrTicker(ticker)
	if sec == nil {
		log.Warnf("couldn't find publisher for %s", ticker)
		if deets := yahoo.GetDetails(ticker); deets != nil && deets.ShortName != "" {
			// Make a new publisher on the fly if we're not already tracking it
			sec := eps.NewStock(deets.ShortName, ticker, "Yahoo")
			pub = eps.NewPublisher(yahoo.GetPrice, ticker, "Yahoo", true, 30)
			pub.UseMarketHours = true
			lookupService.Register(sec, pub)
		} else {
			cb.SendMessage("Could not find details for ticker "+ticker, "", false)
			return
		}
	}

	// Todo: make this a case and extract to funcs
	operation := parts[0]

	if operation == "sub" {
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

	if operation == "get" {
		// Small delay, todo: get rid
		if pub.CurrentCandle == nil {
			time.Sleep(2 * time.Second)
		}
		if pub.CurrentCandle == nil {
			return
		}
		cb.SendMessage(pub.CurrentCandle.String(), "", false)
		cb.SendMessage(pub.StreakSummary(), "", false)
		return
	}

	if operation == "trade" {
		cdl := pub.CurrentCandle
		fee := cdl.Current * .01
		str := "%s -- $%.2f -- Fee: $%.2f -- 2%% Gain: $%.2f ($%.2f)"
		str = fmt.Sprintf(str, cdl.Ticker, cdl.Current, fee, fee*2, fee*2+cdl.Current)
		cb.SendMessage(str, "", false)
	}

	if operation == "chart" || operation == "graph" {
		if queue := queueService.FindByTicker(ticker); queue != nil {
			graph := priceTracking.QueueToGraph(*queue)
			buffer := bytes.NewBuffer([]byte{})
			err := graph.Render(chart.PNG, buffer)
			if err != nil {
				println(err)
				return
			}
			file := &discordgo.File{Name: ticker + ".png", Reader: buffer}
			msg := &discordgo.MessageSend{
				Files: []*discordgo.File{file},
			}
			cb.ds.ChannelMessageSendComplex(cb.channelId, msg)
		} else {
			cb.SendMessage("Couldn't find ticker "+ticker, "", false)
		}
	}
}
