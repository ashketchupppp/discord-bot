package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashketchupppp/discord-bot/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/mattn/go-shellwords"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "token", "", "Discord bot token")
}

func main() {
	flag.Parse()

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
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

func GetCommand(cmd string) (bot.Command, error) {
	args, err := shellwords.Parse(cmd)
	if err != nil {
		return nil, err
	} else {
		if cmd, ok := bot.Commands[args[0]]; ok {
			fs := cmd.FlagSet()
			err := fs.Parse(args[1:])
			if err != nil {
				return nil, err
			}
			return cmd, nil
		}
	}
	return nil, errors.New("unable to find command")
}

func NewMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Content[0] != '/' {
		return
	}
	m.Content = m.Content[1:]
	cmd, err := GetCommand(m.Content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	err = cmd.Run(s, m.Message)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
	}
}
