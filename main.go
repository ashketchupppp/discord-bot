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
	configPath    string
	defaultConfig = &DiscordBot{
		Token: DefaultBotToken,
		MongoDatabase: &MongoDB{
			DBName:          DefaultMongoDBName,
			QuoteCollection: DefaultMongoDBQuoteCollection,
			ConnStr:         DefaultMongoDBConnStr,
		},
		Settings: map[string]string{
			QuoteChannelSettingName: DefaultQuoteChannel,
			LeaveChannelSettingName: DefaultLeaveMessageChannel,
		},
		EnabledFeatures: []string{
			LeaveMessageFeatureName,
		},
		EnabledCommands: []string{
			HelpCmdName,
			AddQuoteCmdName,
			GetQuoteCmdName,
		},
		CommandPrefix: DefaultCommandPrefix,
	}
)

func init() {
	flag.StringVar(&configPath, ConfigPathFlagName, DefaultConfigPath, "Path to the configuration file.")
}

func main() {
	flag.Parse()
	// look for configuration file and read it
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Unable to find the config file at '", configPath, "'. Creating a new one in '", DefaultConfigPath, "'")
		defaultConfigStr, _ := json.Marshal(defaultConfig)
		e := ioutil.WriteFile(DefaultConfigPath, defaultConfigStr, 0)
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
