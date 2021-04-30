package main

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type HelpCommand struct{ name string }

func (cmd *HelpCommand) Name() string { return cmd.name }
func (cmd *HelpCommand) FlagSet() flag.FlagSet {
	fs := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	return *fs
}
func (cmd *HelpCommand) Help() string {
	bot := GetDiscordBot()
	var helpMessage string
	var cmdNum int
	helpMessage += "Here is a list of available commands:\n```"
	cmdNum = 1
	for _, commandName := range bot.EnabledCommands {
		command, _ := bot.GetCommand(commandName)
		if command.Name() != cmd.Name() {
			helpMessage += fmt.Sprint(cmdNum, ". ", command.Help(), "\n")
			cmdNum += 1
		}
	}
	helpMessage += "```"
	return helpMessage
}
func (cmd *HelpCommand) Validate() error {
	return nil
}

func (cmd *HelpCommand) Run(s *discordgo.Session, m *discordgo.Message) error {
	bot := GetDiscordBot()
	bot.SendMessage(cmd.Help(), m.ChannelID, s)
	return nil
}
