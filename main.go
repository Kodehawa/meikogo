package main

import (
	"os"
	"os/signal"
	"syscall"
    "io/ioutil"
	"github.com/bwmarrin/discordgo"
	"encoding/json"
	"log"
	"github.com/go-redis/redis"
)

type Config struct {
	Token string
	OwnerId string
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

var RedisClient *redis.Client
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

	StartRedisClient()

	log.Printf("Starting up Meiko...")

	registerCommands()
	log.Printf("Registered %d commands", len(cmds))

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageWait)

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening up a Websocket connection", err)
		return
	}

	err = dg.UpdateStatus(0, "O-Oh, hi there!")
	if err != nil {
		log.Println(err)
	}

	CheckWaiters()
	anilistTokenUpdate()

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func registerCommands() {
	log.Printf("Registering commands...")
	//Anime commands
	registerCommand("anime", anime())
	registerCommand("character", character())
	//Image commands
	registerCommand("catgirl", catgirl())
	registerCommand("cat", cat())
	//Info commands
	registerCommand("help", help())
	registerCommand("serverinfo", serverinfo())
	registerCommand("userinfo", userinfo())
	registerCommand("ping", ping())
	//Game commands
	registerCommand("trivia", trivia())
	//Config commands
	registerCommand("setprefix", setPrefix())
	//Owner commands
	registerCommand("eval", eval())
}

func registerCommand(name string, cmd Command) {
	cmds[name] = cmd
}

func StartRedisClient() {
	log.Println("Opening Redis Connection...")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
}