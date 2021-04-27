package bot

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	FlagSet() flag.FlagSet
	Run(s *discordgo.Session, m *discordgo.Message) error
	Name() string
}

var (
	Commands map[string]Command
)

func init() {
	Commands = make(map[string]Command)
	addCmd := &AddCommand{name: "Add"}
	Commands[addCmd.Name()] = addCmd
}

type AddCommand struct {
	name string
	One  int
	Two  int
}

func (cmd *AddCommand) Name() string { return cmd.name }

func (cmd *AddCommand) FlagSet() flag.FlagSet {
	fs := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	fs.IntVar(&cmd.One, "one", 0, "")
	fs.IntVar(&cmd.Two, "two", 0, "")
	return *fs
}

func (cmd *AddCommand) Run(s *discordgo.Session, m *discordgo.Message) error {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprint(cmd.One+cmd.Two))
	return nil
}
