package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
    "io/ioutil"
	"github.com/bwmarrin/discordgo"
	"encoding/json"
)

type Config struct {
	Token string
	OwnerId int
	Prefix string
}


var cmds = make(map[string]Command)

var prefix = ""

func main() {
	plan, _ := ioutil.ReadFile("./assets/config.json")
	var data Config
	err := json.Unmarshal(plan, &data)

	Token := data.Token
	prefix = data.Prefix

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	registerCommands()

	dg.AddHandler(messageCreate)


	err = dg.Open()
	if err != nil {
		fmt.Println("boom", err)
		return
	}

	err = dg.UpdateStatus(0,"O-Oh, hi there!")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func registerCommands() {
	registerCommand(ping().Name, ping())
	registerCommand(anime().Name, anime())
	registerCommand(catgirl().Name, catgirl())
}

func registerCommand(name string, cmd Command) {
	cmds[name] = cmd
}