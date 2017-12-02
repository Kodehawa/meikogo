package main

import (
	"github.com/bwmarrin/discordgo"
	"bytes"
	"strings"
	"fmt"
)

type Option struct {
	Help 		string
	Execute 	OptionFunc
	Description string
}

type OptionFunc func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string)

var options = make(map[string]Option)

func conf() (Command) {
	registerOptions()

	return Command{
		Name: "conf",
		Description: "Sets the guild configurations",
		Category: "config",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			args := *split

			if len(args) < 2 {
				s.ChannelMessageSend(message.ChannelID, "H-Huh, you specified an invalid option...")
				return
			}

			var buffer bytes.Buffer
			for i := 0; i < len(args); i++ {
				part := args[i]
				fmt.Println(part)
				if buffer.Len() > 0 {
					buffer.WriteString(":")
				}

				buffer.WriteString(part)

				if _, ok := options[buffer.String()]; ok {
					contentRune := []rune(*content)
					newContent := strings.Trim(string(contentRune[buffer.Len():]), " ")
					newSplit := strings.Split(newContent, " ")
					options[buffer.String()].Execute(s, message, &newContent, &newSplit)
					break
				}
			}
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "Bot Configuration",
				"**This command allows you to check the configuration of the bot. " +
					"There are a handful of different options here that will allow you to personalize your bot a bit more.**", 0xFFB6C1))
		},
	}
}

func registerOptions() {
	options["prefix:set"] = Option {
		Description: "Sets the prefix of this server",
		Help: "",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			args := *split

			perms, err := s.UserChannelPermissions(message.Author.ID, message.ChannelID)
			if err != nil {
				return
			}

			if perms & discordgo.PermissionManageServer > 0  {
				if len(args) < 1 {
					s.ChannelMessageSend(message.ChannelID, "Cannot set a new prefix if you don't tell me the new one!")
				} else {
					prefix := args[0]

					if len(prefix) == 0 {
						s.ChannelMessageSend(message.ChannelID, "Cannot set an empty prefix...")
						return
					}

					id, err := GetGuildId(s, message.ChannelID)
					if err != nil {
						s.ChannelMessageSend(message.ChannelID, "Error while getting the guild this message was sent in...")
						return
					}

					guildData, err := GetGuildData(id)
					if err != nil {
						s.ChannelMessageSend(message.ChannelID, "Error while getting the guild config...")
						return
					}

					guildData.Prefix = prefix
					err = SaveGuildData(id, guildData)

					if err != nil {
						s.ChannelMessageSend(message.ChannelID, "Error while saving the prefix...")
						return
					}

					s.ChannelMessageSend(message.ChannelID, ":white_check_mark: Prefix set successfully to: `" + prefix + "`")
				}
			} else {
				s.ChannelMessageSend(message.ChannelID, ":x: You don't have enough permissions to change the server prefix...")
			}
		},
	}

	options["prefix:reset"] = Option {
		Description: "Resets the server prefix",
		Help: "",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			perms, err := s.UserChannelPermissions(message.Author.ID, message.ChannelID)
			if err != nil {
				return
			}

			if perms & discordgo.PermissionManageServer > 0  {
				id, err := GetGuildId(s, message.ChannelID)
				if err != nil {
					s.ChannelMessageSend(message.ChannelID, "Error while getting the guild this message was sent in...")
					return
				}

				guildData, err := GetGuildData(id)
				if err != nil {
					s.ChannelMessageSend(message.ChannelID, "Error while getting the guild config...")
					return
				}

				guildData.Prefix = ""
				err = SaveGuildData(id, guildData)

				if err != nil {
					s.ChannelMessageSend(message.ChannelID, "Error while saving the prefix...")
					return
				}

				s.ChannelMessageSend(message.ChannelID, ":white_check_mark: Successfully reset prefix!")
			}
		},
	}
}