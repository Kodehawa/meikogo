package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"fmt"
	"time"
	"math/rand"
	"bytes"
	"strconv"
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
					Name: fmt.Sprintf("%s Commands [%d]", strings.Title(k), len(v)),
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
					Text: fmt.Sprintf("Total Commands -> %d", len(cmds)),
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
			channel, err := s.State.Channel(message.ChannelID)
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, ":x: Error while retrieving channel...")
				return
			}

			guild, err := s.State.Guild(channel.GuildID)
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, ":x: Error while retrieving guild...")
				return
			}

			owner, err := s.State.Member(guild.ID, guild.OwnerID)
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, ":x: Error while retrieving owner...")
				return
			}

			roles := guild.Roles
			buffer := bytes.Buffer{}
			for i := 0; i < len(roles); i++ {
				role := roles[i]
				if role.Name != "@everyone" {
					buffer.WriteString(role.Name)
				}

				if i < len(roles) - 1 {
					buffer.WriteString(", ")
				}
			}

			rolesWhole := buffer.String()

			if len(rolesWhole) > 900 {
				split := []rune(rolesWhole)
				rolesWhole = string(split[0:900]) + "..."
			}

			iGID, _ := strconv.ParseInt(guild.ID, 10, 64)
			_time := getCreationTime(iGID)

			s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed {
				Author: &discordgo.MessageEmbedAuthor {
					IconURL: discordgo.EndpointGuildIcon(guild.ID, guild.Icon),
					Name: "Guild Information",
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail {
					URL: discordgo.EndpointGuildIcon(guild.ID, guild.Icon),
				},
				Description: fmt.Sprintf("**Information for %s**\n", guild.Name),
				Fields: []*discordgo.MessageEmbedField {
					{ Name: "ID", Value: guild.ID, Inline: false },
					{ Name: "Channels", Value: fmt.Sprintf("%d", len(guild.Channels)), Inline: true },
					{ Name: "Users", Value: fmt.Sprintf("%d", guild.MemberCount), Inline: true },
					{ Name: "Creation Date", Value: fmt.Sprintf("%d-%d-%d", _time.Year(), _time.Month(), _time.Day()), Inline: true },
					{ Name: "Region", Value: guild.Region, Inline: true },
					{ Name: "Owner", Value: owner.User.Username + "#" + owner.User.Discriminator, Inline: true },
					{ Name: fmt.Sprintf("Roles [%d]", len(guild.Roles)), Value: rolesWhole, Inline: false },
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Requested by " + message.Author.Username + "#" + message.Author.Discriminator,
					IconURL: message.Author.AvatarURL("128"),
				},
				Color: 0x0fa8a5,
			})
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "Server Info",
				"**Shows detailed server information, such as the ID, number of members, owner, roles, etc.**", 0xFFB6C1))
		},
	}
}

func userinfo() (Command) {
	return Command{
		Name: "userinfo",
		Description: "Checks the info of a given user",
		Category: "info",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			mentioned := message.Mentions

			if len(mentioned) == 0 {
				return
			}

			channel, err := s.State.Channel(message.ChannelID)
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, ":x: Error while retrieving channel...")
				return
			}

			user := mentioned[0]
			member, err := s.State.Member(channel.GuildID, user.ID)
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, ":x: Error while retrieving member...")
				return
			}

			roles := member.Roles
			buffer := bytes.Buffer{}
			for i := 0; i < len(roles); i++ {
				role := roles[i]
				if role != "@everyone" {
					buffer.WriteString(role)
				}

				if i < len(roles) - 1 {
					buffer.WriteString(", ")
				}
			}

			rolesWhole := buffer.String()

			if len(rolesWhole) > 900 {
				split := []rune(rolesWhole)
				rolesWhole = string(split[0:900]) + "..."
			}

			iUID, _ := strconv.ParseInt(user.ID, 10, 64)
			_time := getCreationTime(iUID)

			s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{
					IconURL: user.AvatarURL("128"),
					Name: "User Information",
				},
				Description: "**User Information for " + user.Username + "#" + user.Discriminator + "*",
				Fields: []*discordgo.MessageEmbedField{
					{ Name: "ID", Value: user.ID, Inline: false },
					{ Name: "Join Date", Value: member.JoinedAt, Inline: true },
					{ Name: "Creation Date", Value: fmt.Sprintf("%d-%d-%d", _time.Year(), _time.Month(), _time.Day()), Inline: true },
					{ Name: fmt.Sprintf("	Roles [%d]", len(member.Roles)), Value: rolesWhole, Inline: false },
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: user.AvatarURL("256"),
				},
			})

		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "User Info",
				"**Shows user information, such as the ID, discriminator, roles, etc.**", 0xFFB6C1))
		},
	}
}

func about() (Command) {
	return Command {
		Name: "about",
		Description: "Shows general bot information",
		Category: "info",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed {
				Author: &discordgo.MessageEmbedAuthor{
					Name: "About Meiko",
					IconURL: s.State.User.AvatarURL("128"),
				},
				Description: "H-Hi there! I'm **Meiko**, a simple and fun Discord bot made to make your stance on Discord a little bit more dynamic!\n" +
					"To get started and see what I can do, type `//help` \u2764\n\n" +
						"*Some of my features include:*\n" +
						"✓ **Extensive Configuration** (Custom prefixes, welcome messages and more!)" +
						"✓ **Anime Commands** (Check out `//anime`!)\n" +
						"✓ **A lot more to discover** (Check out `//help`!)	",
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Thanks you for choosing me!",
				},
				Fields: []*discordgo.MessageEmbedField {
					{"Servers", fmt.Sprintf("%d", len(s.State.Guilds)), true },
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: s.State.User.AvatarURL("128"),
				},
			})
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {

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