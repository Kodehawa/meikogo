package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	//"fmt"
	"fmt"
)

var categories = make(map[string]map[string]Command)

func help() (Command) {
	return Command {
		Name: "help",
		Description: "Helps you to help the help helping help.",
		Category: "info",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {

			commandArguments := *split

			if len(commandArguments) > 0 {
				if _, ok := cmds[commandArguments[0]]; ok {
					var cmd = cmds[commandArguments[0]]
					cmd.Help(s, message)
				} else {
					s.ChannelMessageSend(message.ChannelID, ":x: That command doesn't exist...")
				}
			}

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
					Name: strings.Title(k),
					Value: getCommandsFromMap(v),
				})
			}

			s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed {
				Title: "Command Help",
				Description: "**Meiko's command help.**\n" +
					"If you need help with a command in particular, please run //help <command>\n",
				Fields: embeds,
				Thumbnail: &discordgo.MessageEmbedThumbnail {
					URL: s.State.User.AvatarURL("512"),
				},
				Footer: &discordgo.MessageEmbedFooter {
					IconURL: message.Author.AvatarURL("128"),
					Text: fmt.Sprintf("Commands ran this session: %d | Total commands: %d", sessionCommands, len(cmds)),
				},
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