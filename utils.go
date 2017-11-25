package main

import (
	"net/http"
	"time"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"fmt"
	"io/ioutil"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	bb, err := ioutil.ReadAll(r.Body)
	if err == nil {
		fmt.Println(string(bb))
	} else {
		fmt.Println(err)
	}

	return json.NewDecoder(r.Body).Decode(&target)
}

func getRawJson(url string) ([]byte, error) {
	r, err := httpClient.Get(url)
	if err != nil {
		return []byte{0}, err
	}

	defer r.Body.Close()

	bb, err := ioutil.ReadAll(r.Body)
	if err == nil {
		return bb, nil
	} else {
		return []byte{0}, err
	}
}

func helpEmbed(s *discordgo.Session, message *discordgo.MessageCreate, name string, content string, color int) (*discordgo.MessageEmbed) {
	return &discordgo.MessageEmbed {
		Thumbnail: &discordgo.MessageEmbedThumbnail {
			URL: "https://cdn3.iconfinder.com/data/icons/line/36/question-512.png",
		},
		Author: &discordgo.MessageEmbedAuthor {
			IconURL: s.State.User.AvatarURL("128"),
			Name: name,
		},
		Description: content,
		Color: color,
		Footer: &discordgo.MessageEmbedFooter {
			Text: "For a list of all commands run //help",
			IconURL: message.Author.AvatarURL("128"),
		},
	}
}