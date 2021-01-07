package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

const secretFile = "discord.json"

var discordSession *discordgo.Session

func discordMessage(str string, atAll bool) {
	for discordSession == nil {
		time.Sleep(1 * time.Second)
	}
	if atAll {
		str += " @everyone"
	}
	discordSession.ChannelMessageSend(conf.Discord.ChannelID, str)
}

func clearChannel() {
	msgs, err := discordSession.ChannelMessages(conf.Discord.ChannelID, 100, "", "", "")
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, m := range msgs {
		go discordSession.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func initBot(token string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	discordSession = dg
	if conf.Discord.ClearChannelOnBoot {
		clearChannel()
	}
	println("Connected to Discord server")
	dg.Close()
}
