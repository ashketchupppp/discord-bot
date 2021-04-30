package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

type GetQuoteCommand struct {
	name   string
	userID string
}

func (cmd *GetQuoteCommand) Name() string { return cmd.name }
func (cmd *GetQuoteCommand) FlagSet() flag.FlagSet {
	fs := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	fs.StringVar(&cmd.userID, "user", "", "")

	// Reset the arguments to defaults so we dont have
	// the values from the last time this was called
	cmd.userID = ""
	return *fs
}
func (cmd *GetQuoteCommand) Help() string {
	bot := GetDiscordBot()
	return fmt.Sprint(bot.CommandPrefix, cmd.Name(), ` -user @user`)
}
func (cmd *GetQuoteCommand) Validate() error {
	return nil
}

func (cmd *GetQuoteCommand) Run(s *discordgo.Session, m *discordgo.Message) error {
	db, err := GetDiscordBot().Database()
	if err != nil {
		return err
	}
	if cmd.userID == "" {
		return errors.New("missing -user")
	}
	if !ValidUserId(cmd.userID) {
		return errors.New("invalid user id")
	}
	user, err := s.User(cmd.userID[3 : len(cmd.userID)-1])
	if err != nil {
		return errors.New("unable to find that user")
	}
	quotes, err := db.GetQuotes(cmd.userID)
	if err != nil {
		return err
	}
	if len(quotes) == 0 {
		return errors.New("unable to find any quotes")
	}
	bot := GetDiscordBot()
	r := rand.New(rand.NewSource(time.Now().Unix()))
	quoteNum := r.Int() % len(quotes)
	bot.SendMessage(fmt.Sprint("```", quotes[quoteNum].Quote, "``` from ", user.Mention()), m.ChannelID, s)

	return nil
}
