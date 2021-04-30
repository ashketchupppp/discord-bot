package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

var singleDiscordBot = &DiscordBot{}

var (
	configPath        string
	defaultConfigPath = "./.bot.conf.json"
	defaultConfig     = &DiscordBot{
		Token: "CHANGE ME",
		MongoDatabase: &MongoDB{
			DBName:          "discordbot",
			QuoteCollection: "quotes",
			ConnStr:         "CHANGE ME",
		},
		Settings: map[string]string{
			"quotechannel": "CHANGE ME",
		},
		EnabledCommands: []string{
			"help",
			"addquote",
			"getquote",
		},
		CommandPrefix: "$",
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
	err = singleDiscordBot.Load(file)
	if err != nil {
		panic(err.Error())
	}
	err = singleDiscordBot.Validate()
	if err != nil {
		panic(err.Error())
	}

	// Initialise our discord bot
	err = singleDiscordBot.Initialise()
	if err != nil {
		panic(err)
	}

	err = singleDiscordBot.Connect()
	if err != nil {
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	singleDiscordBot.Disconnect()
}

// Returns the discord bot currently in use
// This is needed for things like discord event handlers which need access to the
// SingleDiscordBot struct, but can't be passed extra function parameters due to the limitation
// discordgo places on discord event handlers.
func GetDiscordBot() *DiscordBot {
	return singleDiscordBot
}
