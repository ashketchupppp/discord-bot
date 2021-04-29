package main

import (
	"flag"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	FlagSet() flag.FlagSet
	Run(s *discordgo.Session, m *discordgo.Message) error
	Name() string
	Help() string
}
