package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"math/rand"
	"fmt"
)

func ratewaifu() (Command) {
	return Command {
		Name: "ratewaifu",
		Category: "fun",
		Description: "Rates your waifu from 1 to 10.",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			waifu := *content
			rand.Seed(time.Now().Unix())

			s.ChannelMessageSend(message.ChannelID, fmt.Sprintf("I rate %s with a **%d/10**! :eyes:", waifu, rand.Intn(10)))
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "Ratewaifu command", "Rates your waifu (from 1 to 10)!", 0xFFB6C1))
		},
	}
}