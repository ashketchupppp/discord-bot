package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashketchupppp/discord-bot/bot"
	"github.com/ashketchupppp/discord-bot/db"
	"github.com/bwmarrin/discordgo"
)

var (
	configPath        string
	defaultConfigPath = "./.bot.conf.json"
	defaultConfig     = &Config{
		Token:     "CHANGE ME",
		DbConnStr: "mongodb://localhost:27017",
	}
)

type Config struct {
	Token     string
	DbConnStr string
}

func init() {
	flag.StringVar(&configPath, "configPath", defaultConfigPath, "Path to the configuration file.")
}

func main() {
	flag.Parse()
	// look for configuration file and read it
	fileBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Print("Unable to find the config file at '", configPath, "'. Creating a new one in '", defaultConfigPath, "'")
		defaultConfigStr, _ := json.Marshal(defaultConfig)
		e := ioutil.WriteFile(defaultConfigPath, defaultConfigStr, 0)
		if e != nil {
			panic(e)
		}
		return
	}

	var conf *Config
	err = json.Unmarshal(fileBytes, &conf)
	if err != nil {
		panic(err.Error())
	}

	// Establish mongodb connection
	db, err := db.NewBotDB(conf.DbConnStr)
	if err != nil {
		panic(err)
	}
	bot.SetDatabase(db)

	// Establish discord bot connection
	session, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		panic(err)
	}

	// Setup discord event handlers
	session.AddHandler(bot.NewMessageHandler)

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
