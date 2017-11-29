package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
	"fmt"
	"regexp"
)

var sessionCommands = 0
var waiters = make(map[string]Waiter)

type Waiter struct {
	Function WaiterFunc
	Timeout  int64
	Author   string
}

type WaiterFunc func(s *discordgo.Session, m *discordgo.MessageCreate) (bool)

var mentionRegex = regexp.MustCompile("^<@!?321466090159079444>$")

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	go func() {

		botData, err := GetBotData()
		if err != nil {
			botData = &BotData{}
		}

		if stringInSlice(m.Author.ID, botData.BlackListedUsers) {
			return
		}

		var MessageContent = m.Content
		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}

		guildPrefix := prefix

		if stringInSlice(channel.GuildID, botData.BlackListedGuilds) {
			return
		}

		guildData, err := GetGuildData(channel.GuildID)
		if err != nil {
			fmt.Printf("Failed to retrieve data from %s\n", channel.GuildID)
			fmt.Println(err)
		}

		if len(guildData.Prefix) > 0 {
			guildPrefix = guildData.Prefix
		}

		if mentionRegex.MatchString(m.Content) {
			s.ChannelMessageSend(m.ChannelID, ":speech_balloon: Hi, my prefix on this server is `" + guildData.Prefix + "` <3")
			return
		}

		if strings.HasPrefix(MessageContent, guildPrefix) {
			var Content = strings.Replace(MessageContent, guildPrefix, "", 1)
			var SplitContent= strings.Split(Content, " ")
			if len(SplitContent) >= 1 {
				var Command = SplitContent[0]
				SplitContent := SplitContent[1:]
				command, commandExists := cmds[Command]

				if commandExists {
					Content = strings.Trim(strings.Replace(Content, command.Name, "", 1), " ")

					if command.Category == "owner" && m.Author.ID != config.OwnerId {
						s.ChannelMessageSend(m.ChannelID, ":octagonal_sign: You have no permission to execute this command!")
						return
					}

					command.Execute(s, m, &Content, &SplitContent)
					sessionCommands++
				}
			}
		}
	}()
}

func messageWait(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, ok := waiters[m.ChannelID]; ok {
		currentWaiter := waiters[m.ChannelID]
		if m.Author.ID == currentWaiter.Author {
			if currentWaiter.Function(s, m) {
				delete(waiters, m.ChannelID)
			}
		}
	}
}

func CheckWaiters() {
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				for k, w := range waiters {
					if w.Timeout > currentTimeMillis() {
						delete(waiters, k)
					}
				}
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
}
