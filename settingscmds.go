package main

import "github.com/bwmarrin/discordgo"

func setPrefix() (Command) {
	return Command{
		Name: "setprefix",
		Description: "Sets the guild prefix",
		Category: "config",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			args := *split

			if len(args) < 1 {
				s.ChannelMessageSend(message.ChannelID, "Cannot set a new prefix if you don't tell me the new one!")
			} else {
				prefix := args[0]

				channel, err := s.State.Channel(message.ChannelID)
				if err != nil {
					s.ChannelMessageSend(message.ChannelID, "Error while getting the channel this message was sent in...")
					return
				}

				guildData, err := GetGuildData(channel.GuildID)
				if err != nil {
					s.ChannelMessageSend(message.ChannelID, "Error while getting the guild config...")
					return
				}

				guildData.Prefix = prefix
				err = SaveGuildData(channel.GuildID, guildData)

				if err != nil {
					s.ChannelMessageSend(message.ChannelID, "Error while saving the prefix...")
					return
				}

				s.ChannelMessageSend(message.ChannelID, ":white_check_mark: Prefix set successfully to: `" + prefix + "`")
			}
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {

		},
	}
}