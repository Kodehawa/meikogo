package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"time"
	"net/url"
	"io/ioutil"
	"log"
	"encoding/json"
	"strings"
	"bytes"
)

type AnimeData struct {
	Id            int       `json:"id"`
	Description   string    `json:"description"`
	AverageScore  int 		`json:"average_score"`
	Duration      int       `json:"duration"`
	EndDate       string    `json:"end_date"`
	Genres        []string  `json:"genres"`
	SeriesType    string    `json:"series_type"`
	Type          string    `json:"type"`
	LargeImageUrl string 	`json:"image_url_lge"`
	SmallImageUrl string 	`json:"image_url_sml"`
	MedImageUrl   string	`json:"image_url_med"`
	StartDate     string    `json:"start_date"`
	EnglishTitle  string    `json:"title_english"`
	JapaneseTitle string    `json:"title_japanese"`
	RomajiTitle   string    `json:"title_romaji"`
	TotalEpisodes int		`json:"total_episodes"`
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

			if len(keys) == 0 {
				s.ChannelMessageSend(message.ChannelID, ":x: There are no animes matching your search...")
				return
			}

			if len(keys) > 1 {
				var buffer bytes.Buffer
				buffer.WriteString("\n")

				for i, v := range keys {
					if i > 4 {
						break
					}

					buffer.WriteString(fmt.Sprintf("**%d**.- *%s*", i + 1, v.EnglishTitle))
					if i < 3 || i < len(keys) {
						buffer.WriteString("\n")
					}
				}

				s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed {
					Title: "Anime Selection. Type the number to continue",
					Description: "\n" + buffer.String(),
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "http://i.imgur.com/apTJEIj.png",
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Information provided by Anilist",
						IconURL: message.Author.AvatarURL("128"),
					},
				})

				converted := make([]interface{}, len(keys))
				for i, v := range keys {
					converted[i] = v
				}

				createSelectionWaiter(converted, message, func(key interface{}){
					anime := key.(AnimeData)
					animeInfo(anime, s, message)
				})
			} else {
				anime := keys[0]
				animeInfo(anime, s, message)
			}

		},
		Help: func(s *discordgo.Session, message *discordgo.MessageCreate) {
			s.ChannelMessageSendEmbed(message.ChannelID, helpEmbed(s, message, "Anime Command", "**Search for your favorite anime!**", 0xa8a5c9))
		},
	}
}

func animeInfo(anime AnimeData, s *discordgo.Session, message *discordgo.MessageCreate) {
	descriptionWhole := anime.Description

	if len(descriptionWhole) > 1200 {
		description := []rune(anime.Description)
		descriptionWhole = string(description[0:1200]) + "..."
	}

	s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed {
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: message.Author.AvatarURL("128"),
			Name: fmt.Sprintf("Information of %s (%s)\n\n", anime.EnglishTitle, anime.JapaneseTitle),
		},
		Description: "\n" +  strings.Replace(descriptionWhole, "<br>", "\n", 10),
		Fields: []*discordgo.MessageEmbedField{
			{ Name: "Score", Value: fmt.Sprintf("%d",anime.AverageScore) + "/100", Inline: true, },
			{ Name: "Type", Value: strings.Title(anime.Type) , Inline: true, },
			{ Name: "Start Date", Value: strings.Split(anime.StartDate, "T")[0] , Inline: true, },
			{ Name: "End Date", Value: strings.Split(anime.EndDate, "T")[0] , Inline: true, },
			{ Name: "Genres", Value: strings.Join(anime.Genres, ", "), Inline: false, },
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: anime.MedImageUrl,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Information provided by Anilist",
			IconURL: "http://i.imgur.com/apTJEIj.png",
		},
	})
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