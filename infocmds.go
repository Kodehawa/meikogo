package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"fmt"
	"time"
	"math/rand"
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

				return
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
				Author: &discordgo.MessageEmbedAuthor{
					IconURL: s.State.User.AvatarURL("128"),
					Name: "Meiko Help",
				},
				Description: "**Meiko's command help.**\n" +
					"For extended command usage please run *//help <command>*\n",
				Fields: embeds,
				Footer: &discordgo.MessageEmbedFooter {
					IconURL: message.Author.AvatarURL("128"),
					Text: fmt.Sprintf("Commands ran this session: %d | Total commands: %d", sessionCommands, len(cmds)),
				},
				Color: 0x37b75b,
			})
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "Help Command", "**Well, this command.**", 0xFFB6C1))
		},
	}
}

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

func serverinfo() (Command) {
	return Command {
		Name: "serverinfo",
		Description: "Shows the information of this server",
		Category: "info",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {

		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "Server Info",
				"**Shows detailed server information, such as the ID, number of members, owner, roles, etc.**", 0xFFB6C1))
		},
	}
}

func getCommandsFromMap(m map[string]Command) (string) {
	var parts []string

	for _, v := range m {
		parts = append(parts, "`" + v.Name + "`")
	}

	return strings.Join(parts, ", ")
}