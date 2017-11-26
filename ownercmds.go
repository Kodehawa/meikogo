package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	eval2 "github.com/apaxa-go/eval"
	"reflect"
)

func eval() (Command) {
	return Command{
		Name: "eval",
		Description: "Evaluates arbitrary code.",
		Category: "owner",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			expression := *content
			expr, err := eval2.ParseString(expression, "")
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, fmt.Sprintf(":x: Error while evaluating this expression! (%s)", err))
				return
			}

			a := eval2.Args {
				"message": 		 eval2.MakeDataRegular(reflect.ValueOf(message)),
				"session": 		 eval2.MakeDataRegular(reflect.ValueOf(s)),
				"fmt.Sprint":    eval2.MakeDataRegularInterface(fmt.Sprint),
			}

			r, err := expr.EvalToInterface(a)
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, fmt.Sprintf(":x: Error while evaluating this expression! (%s)", err))
				return
			}

			s.ChannelMessageSend(message.ChannelID, fmt.Sprintf(":white_check_mark: Evaluated with success and returned: %v", r))
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "Eval Command", "Evaluates arbitrary code", 0xFFB6C1))
		},
	}
}