package main

import (
	"os"
	"os/signal"
	"syscall"
    "io/ioutil"
	"github.com/bwmarrin/discordgo"
	"encoding/json"
	"log"
)

type Config struct {
	Token string
	OwnerId int
	Prefix string
	AnilistToken string `json:"anilist_token"`
	AnilistSecret string `json:"anilist_secret"`
}

var cmds = make(map[string]Command)
var prefix = ""
var config Config

func main() {
	plan, _ := ioutil.ReadFile("./assets/config.json")
	err := json.Unmarshal(plan, &config)
	if err != nil {
		log.Fatal("Error parsing config json file!")
		return
	}

	Token := config.Token
	prefix = config.Prefix

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}

	log.Printf("Starting up Meiko...")

	registerCommands()
	log.Printf("Registered %d commands", len(cmds))

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening up a Websocket connection", err)
		return
	}

	err = dg.UpdateStatus(0,"O-Oh, hi there!")
	if err != nil {
		log.Println(err)
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func registerCommands() {
	registerCommand(ping().Name, ping())
	registerCommand(anime().Name, anime())
	registerCommand(catgirl().Name, catgirl())
	registerCommand(help().Name, help())
}

func registerCommand(name string, cmd Command) {
	cmds[name] = cmd
}