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
	AnilistClient string   `json:"anilist_key"`
	AnilistSecret string   `json:"anilist_secret"`
	WeatherToken string    `json:"weatherAppId"`
	DBotsOrgToken string
	DBotsToken string
}

type Command struct {
	Name string
	Description string
	Category string
	Execute HandlerFunc
	Help HelpFunc
}

type HandlerFunc func(s *discordgo.Session, message *discordgo.MessageCreate, content *string, split *[]string)
type HelpFunc func(s *discordgo.Session, message *discordgo.MessageCreate)

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

	err = dg.UpdateStatus(0, "O-Oh, hi there!")
	if err != nil {
		log.Println(err)
	}

	anilistTokenUpdate()
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func registerCommands() {
	registerCommand("ping", ping())
	registerCommand("anime", anime())
	registerCommand("catgirl", catgirl())
	registerCommand("help", help())
}

func registerCommand(name string, cmd Command) {
	cmds[name] = cmd
}