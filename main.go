package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	configPath        string
	defaultConfigPath = "./.bot.conf.json"
	defaultConfig     = &DiscordBot{
		Token:     "CHANGE ME",
		DbConnStr: "CHANGE ME",
	}
)

func init() {
	flag.StringVar(&configPath, "configPath", defaultConfigPath, "Path to the configuration file.")
}

func main() {
	flag.Parse()
	// look for configuration file and read it
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Unable to find the config file at '", configPath, "'. Creating a new one in '", defaultConfigPath, "'")
		defaultConfigStr, _ := json.Marshal(defaultConfig)
		e := ioutil.WriteFile(defaultConfigPath, defaultConfigStr, 0)
		if e != nil {
			panic(e)
		}
		return
	}

	// We have a config file, read it and validate the discordBot is setup correctly
	var discordBot *DiscordBot
	err = discordBot.Load(file)
	if err != nil {
		panic(err.Error())
	}
	err = discordBot.Validate()
	if err != nil {
		panic(err.Error())
	}

	// Establish mongodb connection
	db, err := NewBotDB(discordBot.DbConnStr)
	if err != nil {
		panic(err)
	}
	SetDatabase(db)

	// Establish discord bot connection
	session, err := discordgo.New("Bot " + discordBot.Token)
	if err != nil {
		panic(err)
	}

	// Setup discord event handlers
	session.AddHandler(NewMessageHandler)

	err = session.Open()
	if err != nil {
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}
