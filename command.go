package main

import "github.com/bwmarrin/discordgo"

type Command struct {
	Name string
	Description string
	Category string
	Execute func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string)
	Help func(s *discordgo.Session, message *discordgo.MessageCreate)
}