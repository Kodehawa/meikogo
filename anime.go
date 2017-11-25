package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"time"
	"net/url"
	"io/ioutil"
	"log"
	"encoding/json"
)

type AnimeData struct {
	AverageScore  string    `json:"average_score"`
	Description   string    `json:"description"`
	Duration      int       `json:"duration"`
	EndDate       string    `json:"end_date"`
	Genres        []string  `json:"genres"`
	Id            int       `json:"id"`
	LargeImageUrl string    `json:"image_url_lge"`
	SmallImageUrl string    `json:"image_url_sml"`
	SeriesType    string    `json:"series_type"`
	StartDate     string    `json:"start_date"`
	EnglishTitle  string    `json:"title_english"`
	JapaneseTitle string    `json:"title_japanese"`
	RomajiTitle   string    `json:"title_romaji"`
	TotalEpisodes string    `json:"total_episodes"`
}

type AnilistData struct {
	AuthToken string    `json:"access_token"`
}

var AniListData = &AnilistData{}

func anime() (Command) {
	return Command {
		Name: "anime",
		Description: "Search for your favorite anime!",
		Category: "anime",
		Execute: func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string) {
			splitContent := *split
			if len(splitContent) < 1 {
				s.ChannelMessageSend(message.ChannelID, ":x: You need to specify the name of the anime!")
				return
			}

			url1, err := url.Parse(fmt.Sprintf("https://anilist.co/api/anime/search/%s?access_token=%s", *content, AniListData.AuthToken))
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, ":sob: Uh... an error happening while retrieving this anime :<")
				return
			}

			keys := make([]AnimeData, 0)
			response, err := getRawJson(url1.String())
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, ":sob: Uh... an error happening while retrieving this anime data :<")
				return
			}

			json.Unmarshal(response, &keys)
			s.ChannelMessageSend(message.ChannelID, keys[0].Description)
		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "Anime Command", "**Search for your favorite anime!**", 0xa8a5c9))
		},
	}
}

func anilistTokenUpdate () {
	log.Println("Starting Anilist authentication token task...")
	updateToken()

	ticker := time.NewTicker(30 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				updateToken()
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func updateToken() {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", config.AnilistClient)
	form.Add("client_secret", config.AnilistSecret)
	r, err := httpClient.PostForm("https://anilist.co/api/auth/access_token", form)
	r.Close = true

	if err != nil {
		log.Fatal("Cannot update auth token!", err)
		return
	}

	bb, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal("Cannot read response from Anilist!", err)
		return
	}

	json.Unmarshal(bb, AniListData)
	log.Println("Updated AniList authentication token")
}