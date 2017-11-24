package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	go func() {
		var MessageContent = m.Content
		if strings.HasPrefix(MessageContent, prefix) {
			var Content = strings.Replace(MessageContent, prefix, "", 1)
			var SplitContent= strings.Split(Content, " ")
			if len(SplitContent) >= 1 {
				var Command = SplitContent[0]
				command, commandExists := cmds[Command]

				if commandExists {
					command.Execute(s, m, &Content, &SplitContent)
				}
			}
		}
	}()
}
