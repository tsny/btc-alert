package main

import (
	"testing"
)

func Test_GetTopGainers(t *testing.T) {
	cryptoBot.GetTopGainers()
}

func TestTTS(t *testing.T) {
	// cryptoBot.SendMessage("test", "tsny", true)
	// cryptoBot.ds.ChannelMessageSendTTS(conf.Discord.ChannelID, "test")

	// msg := discordgo.MessageSend{
	// 	Content: "test @everyone",
	// 	TTS:     true,
	// 	AllowedMentions: &discordgo.MessageAllowedMentions{
	// 		Users: []string{"84090395092353024"},
	// 		Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeEveryone},
	// 	},
	// }
	// cryptoBot.ds.ChannelMessageSendComplex(conf.Discord.ChannelID, &msg)
}

// func TestNotification(t *testing.T) {
// 	i := interval{5, 0, 3, 0, time.Now()}
// 	bannerText := "$33945.07 --> $33867.83 | Change: $-77.24 | Percent: -0.228%"
// 	hdr := sf("%d Minutes Passed | %.2f%%", i.MaxOccurences, i.PercentThreshold)
// 	notif(hdr, bannerText, "assets/warning.png")
// }
