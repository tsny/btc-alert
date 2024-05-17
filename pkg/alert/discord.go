package alert

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

// CryptoBot is a service that communicates with discord and holds onto alerts
// that are created for discord users
type CryptoBot struct {
	session *discordgo.Session
}

// todo: pass in the channel id, don't ref config all the time?
// todo: return err
func NewBot(token string) *CryptoBot {
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
	cb := &CryptoBot{session: session}
	session.AddHandler(cb.OnNewMessage)
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

	case "boot":
		str := fmt.Sprintf("boot time: %v (%v ago)", bootTime.Local().Format(time.DateTime), time.Since(bootTime))
		_, _ = cb.SendUserMessage(m.Author, str)

	case "whois", "get":
		cb.handleGet(ticker, m)

	case "target":
		if len(parts) <= 2 {
			_, _ = cb.SendUserMessage(m.Author, "Invalid; usage: target btc 60000")
			return
		}
		pub, ok := FindPublisher(Publishers, ticker)
		if !ok {
			_, _ = cb.SendUserMessage(m.Author, "ticker %v not found", ticker)
			return
		}
		target, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			_, _ = cb.SendUserMessage(m.Author, "coun't parse num %v: %v", parts[2], err)
			return
		}
		NewChangeListener(pub).RegisterTargetTracker(m.Author.ID, target)
		_, _ = cb.SendUserMessage(m.Author, "Will alert you when `%v` price (`%v`) moves past `%v`",
			ticker, pub.Candle.Price, target)

	case "track":
		if len(parts) <= 2 {
			_, _ = cb.SendUserMessage(m.Author, "Invalid; usage: track btc 500")
			return
		}
		pub, ok := FindPublisher(Publishers, ticker)
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
			cl.RegisterPriceMovementListener(m.Author.ID, chgAmount)
		}
		_, _ = cb.SendUserMessage(m.Author, "Will alert you when `%v` price (`%v`) changes by `%v%v`",
			ticker, pub.Candle.Price, chgAmount, suffix)

	default:
		_, _ = cb.SendUserMessage(m.Author, "What?")
	}
}

// SendMessage sends a discord message with an optional mention
// TODO: Could change these into options for UserID and TTS
func (cb *CryptoBot) SendMessage(str string, userID string) (*discordgo.Message, error) {
	log.Infof("Sending message to %v", userID)
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

func (cb *CryptoBot) handleGet(ticker string, m *discordgo.MessageCreate) {
	pub, ok := FindPublisher(Publishers, ticker)
	if !ok {
		_, _ = cb.session.ChannelMessageSend(m.ChannelID, "Who?")
		return
	}
	log.Infof("%v", pub.String())
	_, _ = cb.session.ChannelMessageSend(m.ChannelID, pub.String())
}
