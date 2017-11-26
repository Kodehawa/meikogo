package main

import (
	"github.com/bwmarrin/discordgo"
	"encoding/json"
	"math/rand"
	"html"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)


type Trivia struct {
	Results []Result  `json:"results"`
}

type Result struct {
	Category 		 string	  `json:"category"`
	Type 	 		 string	  `json:"type"`
	Difficulty  	 string	  `json:"difficulty"`
	Question 		 string	  `json:"question"`
	CorrectAnswer 	 string	  `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

func trivia() (Command) {
	return Command{
		Name: "trivia",
		Description: "Starts a game of trivia!",
		Category: "game",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			t, err := setup()
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, ":x: Unexpected error while setting up a game of trivia...")
				return
			}

			if len(t.Results) == 0 {
				s.ChannelMessageSend(message.ChannelID, ":x: Unexpected error while setting up a game of trivia...")
				return
			}

			result := t.Results[0]
			answers := result.IncorrectAnswers
			answers = append(answers, result.CorrectAnswer)

			//Shuffle the list...
			//Also takes care of giving me a human-readable output.
			for i := len(answers) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				answers[i], answers[j] = html.UnescapeString(answers[j]), html.UnescapeString(answers[i])
			}

			correctNumber := 0
			buffer := bytes.Buffer{}
			for i:= 0; i < len(answers); i++ {

				if answers[i] == result.CorrectAnswer {
					correctNumber = i + 1
				}

				buffer.WriteString(fmt.Sprintf("**%d.-** %s\n", i + 1 , answers[i]))
			}

			s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
				Title: "Trivia Game",
				Description: "Question: **" + html.UnescapeString(result.Question) + "**" + "\n\n" + buffer.String(),
				Footer: &discordgo.MessageEmbedFooter{
					Text: "You have 60 seconds and 3 attempts to answer",
					IconURL: message.Author.AvatarURL("128"),
				},
				Color: 0x8cd66f,
			})

			attempts := 0

			CreateWaiter(message.ChannelID, message.Author.ID, func(s *discordgo.Session, m *discordgo.MessageCreate) bool {
				if strings.ToLower(m.Content) == "end" {
					s.ChannelMessageSend(message.ChannelID, ":ok_hand: Ended trivia.")
					return true
				}

				i, err := strconv.ParseInt(m.Content, 10, 32)
				if err != nil {
					if strings.ToLower(html.UnescapeString(result.CorrectAnswer)) == strings.ToLower(m.Content) {
						s.ChannelMessageSend(message.ChannelID, ":tada: Correct answer!")
						return true
					}

					attempts++
					if attempts > 2 {
						s.ChannelMessageSend(message.ChannelID, ":sob: Already used all attempts, correct answer was *" + html.UnescapeString(result.CorrectAnswer) + "*")
						return true
					} else {
						s.ChannelMessageSend(message.ChannelID, fmt.Sprintf(":warning: **Incorrect answer!** (Attempts remaning: %d)", 3 - attempts))
						return false
					}
				} else {
					if i > 5 {
						return false
					}

					if correctNumber == int(i) {
						userData, err := GetUserData(message.Author.ID)
						if err == nil {
							userData.IncrementGamesWon()
							SaveUserData(message.Author.ID, userData)
						}
						s.ChannelMessageSend(message.ChannelID, fmt.Sprintf(":tada: Correct answer! Total Games you've won: %d", userData.GamesWon))
						return true
					}
				}

				return false //Missing return at end of function even though it never reaches here? :thinking:
			})

		}, Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID,helpEmbed(s, message, "Trivia Game", "**Starts a game of Trivia!**", 0xFFB6C1))
		},
	}
}

func setup() (t *Trivia, err error) {
	res, err := getRawJson("https://opentdb.com/api.php?amount=1")
	trivia := &Trivia{}

	if err != nil {
		return nil, err
	}

	json.Unmarshal(res, &trivia)

	return trivia, nil
}