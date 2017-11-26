package main

import "github.com/bwmarrin/discordgo"

func trivia() (Command) {
	return Command{
		Name: "trivia",
		Description: "Starts a game of trivia!",
		Category: "game",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {

		}, Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {

		},
	}
}