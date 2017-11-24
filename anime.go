package main

import "github.com/bwmarrin/discordgo"

func anime() (Command) {
	return Command{
		Name: "anime",
		Description: "Search for your favorite anime!",
		Category: "anime",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			splitContent := *split

			if len(splitContent) < 2 {

			}
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {

		},
	}
}