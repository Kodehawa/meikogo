package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
    "math/rand"
    "fmt"
)

func ping() (Command) {
	pingQuotes := [4]string{"Pong", "D-Did I do well?", "Oh oh, look at what I can do!", "Woaaah~"}

	return Command {
		 Name: "ping",
		 Description: "Plays ping-pong with the user",
		 Category: "info",
		 Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
		 	var start = time.Now().UnixNano() / 1000000
		 	s.ChannelTyping(message.ChannelID)
		 	var end = time.Now().UnixNano() / 1000000

		 	rand.Seed(time.Now().Unix())

			 //Sends the pong :tm:
		 	s.ChannelMessageSend(message.ChannelID, fmt.Sprintf(":mega: *%s.* I took **%d ms** to get back to you! :heart:", pingQuotes[rand.Intn(4)], end - start))
		 },
		 Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			 s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "Ping Command", "**Displays the bot's ping.**", 0xFFB6C1))
		 },
	 }
}