package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"fmt"
)

type GuildData struct {
	Prefix 		    string `json:"prefix"`
	WelcomeMessage  string `json:"welcome_message"`
	LeaveMessage    string `json:"leave_message"`
	NoDefaultPrefix bool   `json:"no_default_prefix"`
}

type BotData struct {
	BlackListedUsers  []string 	`json:"blacklisted_users"`
	BlackListedGuilds []string 	`json:"blacklisted_guilds"`
}

type UserData struct {
	GamesWon   int32 	`json:"games_won"`
	Experience int32 	`json:"experience"`
	Waifu 	   string   `json:"waifu"`

}

func GetBotData() (*BotData, error) {
	bd, err := RedisClient.Get("meiko_data").Result()
	if err == redis.Nil {
		SaveBotData(&BotData{})
		data, err := GetBotData()
		if err != nil {
			return nil, err
		}

		return data, nil
	} else if err != nil {
		return nil, err
	}

	botData := &BotData{}
	err = json.Unmarshal([]byte(bd), botData)
	if err != nil {
		return nil, err
	}

	return botData, nil
}

func SaveBotData(data *BotData) error {
	serializedJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = RedisClient.Set("meiko_data", serializedJson, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetGuildData(guildID string) (*GuildData, error) {
	j, err := RedisClient.HGet("guilds", guildID).Result()
	if err == redis.Nil {
		fmt.Println("Cannot find redis data for guild " + guildID + ", saving new one...")
		SaveGuildData(guildID, &GuildData{})
		data, err := GetGuildData(guildID)
		if err != nil {
			return nil, err
		}

		return data, nil
	} else if err != nil {
		return nil, err
	}

	guildData := &GuildData{}
	err = json.Unmarshal([]byte(j), guildData)
	if err != nil {
		return nil, err
	}

	return guildData, nil
}

func SaveGuildData(guildID string, data *GuildData) error {
	serializedJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = RedisClient.HSet("guilds", guildID, serializedJson).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUserData(userID string) (*UserData, error) {
	u, err := RedisClient.HGet("users", userID).Result()
	if err == redis.Nil {
		fmt.Println("Cannot find redis data for user " + userID + ", saving new one...")
		SaveUserData(userID, &UserData{})
		data, err := GetUserData(userID)
		if err != nil {
			return nil, err
		}

		return data, nil
	} else if err != nil {
		return nil, err
	}

	userData := &UserData{}
	err = json.Unmarshal([]byte(u), userData)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func SaveUserData(userId string, data *UserData) error {
	serializedJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = RedisClient.HSet("users", userId, serializedJson).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ud *UserData) IncrementExperience() {
	experience := ud.Experience
	ud.Experience = experience + 1
}

func (ud *UserData) IncrementGamesWon() {
	gamesWon := ud.GamesWon
	ud.GamesWon = gamesWon + 1
}
