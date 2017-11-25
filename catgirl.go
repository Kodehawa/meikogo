package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
)

type Neko struct {
	Url string `json:"neko"`
}

func catgirl() (Command) {
	return Command {
		Name: "catgirl",
		Description: "N-Nya~",
		Category: "image",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			splitContent := *split
			neko := &Neko{}
			if len(splitContent) == 1 {
				if splitContent[0] == "nsfw" || splitContent[0] == "lewd" {
					channel, err := s.Channel(message.ChannelID)
					if err != nil {
						fmt.Println(err)
					} else {
						if channel.NSFW {
							SendNekoImage(s, message, "https://nekos.life/api/lewd/neko", neko)
						} else {
							s.ChannelMessageSend(message.ChannelID, ":x: You lewdie, but you cannot use NSFW commands outside of NSFW channels~")
						}
					}
				} else {
					s.ChannelMessageSend(message.ChannelID, ":x: Incorrrect arguments :<")
				}
			} else {
				SendNekoImage(s, message, "https://nekos.life/api/neko", neko)
			}
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {}, //TODO make command help
	}
}

func SendNekoImage(s *discordgo.Session, message *discordgo.MessageCreate, url string, target *Neko){
	err := getJson(url, target)

	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(message.ChannelID, fmt.Sprintf(":x: Something went wrong while looking for this image... (%s)", err))
		return
	}

	s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed {
		Description: "**N-Nya~**",
		Image: &discordgo.MessageEmbedImage{
			URL: target.Url,
		},
		Color: 0x07beb8,
	})
}