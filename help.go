package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

var categories = make(map[string]map[string]Command)

func help() (Command) {
	return Command {
		Name: "help",
		Description: "Helps you to help the help helping help.",
		Category: "info",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			var embeds []*discordgo.MessageEmbedField

			for name, value := range cmds {
				if categories[value.Category] == nil {
					categories[value.Category] = make(map[string]Command)
				}

				if _, ok := categories[value.Category][name]; !ok {
					categories[value.Category][name] = value
				}
			}

			for k, v := range categories {
				embeds = append(embeds, &discordgo.MessageEmbedField{
					Name: strings.ToUpper(k),
					Value: getCommandsFromMap(v),
				})
			}

			s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
				Title: "Help",
				Fields: embeds,
			})
 		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {},
	}
}

func getCommandsFromMap(m map[string]Command) (string) {
	var parts []string

	for _, v := range m {
		parts = append(parts, "`" + v.Name + "`")
	}

	return strings.Join(parts, ", ")
}