package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type AddQuoteCommand struct {
	name         string
	userID       string
	quote        string
	quoteChannel string
}

func (cmd *AddQuoteCommand) Name() string { return cmd.name }
func (cmd *AddQuoteCommand) FlagSet() flag.FlagSet {
	fs := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	fs.StringVar(&cmd.userID, "user", "", "")
	fs.StringVar(&cmd.quote, "quote", "", "")

	// Reset the arguments to defaults so we dont have
	// the values from the last time this was called
	cmd.userID = ""
	cmd.quote = ""
	cmd.quoteChannel = ""
	return *fs
}
func (cmd *AddQuoteCommand) Help() string {
	bot := GetDiscordBot()
	return fmt.Sprint(bot.CommandPrefix, cmd.Name(), " -user @User -quote \"quote\"")
}
func (cmd *AddQuoteCommand) Validate() error {
	bot := GetDiscordBot()
	val, err := bot.GetSetting("quotechannel")
	if err == nil {
		_, err = strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("the setting quotechannel must be an integer")
		}
		if len(val) != 18 {
			return fmt.Errorf("the setting quotechannel must be of length 18")
		}
	}
	return nil
}

func (cmd *AddQuoteCommand) Run(s *discordgo.Session, m *discordgo.Message) error {
	db, err := GetDiscordBot().Database()
	if err != nil {
		return err
	}
	if cmd.userID == "" {
		return errors.New("missing -user")
	}
	if cmd.quote == "" {
		return errors.New("missing -quote")
	}
	if !ValidUserId(cmd.userID) {
		return errors.New("invalid user id")
	}
	user, err := s.User(cmd.userID[3 : len(cmd.userID)-1])
	if err != nil {
		return errors.New("unable to find that user")
	}
	bot := GetDiscordBot()
	db.AddQuote(Quote{Time: time.Now().Format(time.RFC850), User: cmd.userID, Quote: cmd.quote})
	bot.SendMessage(fmt.Sprint("Added a new quote for ", user.Mention()), m.ChannelID, s)
	if qchannel, err := bot.GetSetting("quotechannel"); err == nil {
		bot.SendMessage(fmt.Sprint("```", cmd.quote, "```", " from ", user.Mention()), qchannel, s)
	}
	return nil
}
