package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"fmt"
)

type GuildData struct {
	Prefix string `json:"prefix"`
}

type BotData struct {
	BlackListedUsers  []string 	`json:"blacklisted_users"`
	BlackListedGuilds []string 	`json:"blacklisted_guilds"`
}

type UserData struct {
	GamesWon int32 	`json:"games_won"`
}

func GetGuildData(guildID string) (*GuildData, error) {
	j, err := RedisClient.Get("guild:" + guildID).Result()
	if err == redis.Nil {
		fmt.Println("Cannot find redis data for " + guildID + ", saving new one...")
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

	err = RedisClient.Set("guild:" + guildID, serializedJson, 0).Err()

	if err != nil {
		return err
	}

	return nil
}