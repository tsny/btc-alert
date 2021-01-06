package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type discordConfig struct {
	Token              string `json:"token"`
	ChannelID          string `json:"channelId"`
	ClearChannelOnBoot bool   `json:"clearOnBoot"`
}

const secretFile = "discord.json"

var discordConf *discordConfig
var discordSession *discordgo.Session

func initToken() string {
	println("enabling discord integration")
	bytes, err := ioutil.ReadFile("discord.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, &discordConf)
	if len(discordConf.Token) < 1 {
		log.Fatalf("Didn't get a key from '%s'", secretFile)
	}
	return discordConf.Token
}

func discordMessage(str string, atAll bool) {
	for discordSession == nil {
		time.Sleep(1 * time.Second)
	}
	if atAll {
		str += " @all"
	}
	discordSession.ChannelMessageSend(discordConf.ChannelID, str)
}

func clearChannel() {
	msgs, err := discordSession.ChannelMessages(discordConf.ChannelID, 100, "", "", "")
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
	if discordConf.ClearChannelOnBoot {
		clearChannel()
	}
	dg.Close()
}
