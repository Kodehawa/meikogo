package main

import "github.com/bwmarrin/discordgo"

func setPrefix() (Command) {
	return Command{
		Name: "setprefix",
		Description: "Sets the guild prefix",
		Category: "config",
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
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {

		},
	}
}