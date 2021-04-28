package bot

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/ashketchupppp/discord-bot/db"
	"github.com/bwmarrin/discordgo"
	"github.com/mattn/go-shellwords"
)

type Command interface {
	FlagSet() flag.FlagSet
	Run(s *discordgo.Session, m *discordgo.Message) error
	Name() string
	Help() string
}

var (
	Commands      map[string]Command
	CommandPrefix byte
	botDatabase   *db.BotDB
)

func init() {
	CommandPrefix = '$'
	Commands = make(map[string]Command)

	helpCmd := &HelpCommand{name: "help"}
	Commands[helpCmd.Name()] = helpCmd
	addQuoteCmd := &AddQuoteCommand{name: "addquote"}
	Commands[addQuoteCmd.Name()] = addQuoteCmd
	getQuoteCommand := &GetQuoteCommand{name: "getquote"}
	Commands[getQuoteCommand.Name()] = getQuoteCommand
}

func getCommand(cmdStr string) (Command, error) {
	// if the command contains user ids, they need quotes around them
	re := regexp.MustCompile(UserIDRegex)
	userIDs := re.FindAll([]byte(cmdStr), 10)
	for id := range userIDs {
		cmdStr = strings.Replace(cmdStr, string(userIDs[id]), fmt.Sprint("\"", string(userIDs[id]), "\""), -1)
	}

	args, err := shellwords.Parse(cmdStr)
	if err != nil {
		return nil, err
	}
	cmd, cmdExists := Commands[args[0]]
	if !cmdExists {
		return nil, errors.New("unable to find command")
	}
	fs := cmd.FlagSet()
	err = fs.Parse(args[1:])
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func sendError(e error, chanID string, s *discordgo.Session) {
	sendMessage(fmt.Sprint("```", e.Error(), "\nfor help, use $help", "```"), chanID, s)
}

func sendMessage(msg, chanID string, s *discordgo.Session) {
	s.ChannelMessageSend(chanID, msg)
}

func NewMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Content[0] != CommandPrefix {
		return
	}
	m.Content = m.Content[1:]
	cmd, err := getCommand(m.Content)
	if err != nil {
		sendError(err, m.ChannelID, s)
		return
	}
	err = cmd.Run(s, m.Message)
	if err != nil {
		sendError(err, m.ChannelID, s)
	}
}

func SetDatabase(db *db.BotDB) {
	botDatabase = db
}

type HelpCommand struct{ name string }

func (cmd *HelpCommand) Name() string { return cmd.name }
func (cmd *HelpCommand) FlagSet() flag.FlagSet {
	fs := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	return *fs
}
func (cmd *HelpCommand) Help() string {
	var helpMessage string
	var cmdNum int
	helpMessage += "Here is a list of available commands:\n```"
	cmdNum = 1
	for _, command := range Commands {
		if command.Name() != cmd.Name() {
			helpMessage += fmt.Sprint(cmdNum, ". ", command.Help(), "\n")
			cmdNum += 1
		}
	}
	helpMessage += "```"
	return helpMessage
}
func (cmd *HelpCommand) Run(s *discordgo.Session, m *discordgo.Message) error {
	sendMessage(cmd.Help(), m.ChannelID, s)
	return nil
}

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
	return fmt.Sprint(string(CommandPrefix), cmd.Name(), " -user @User -quote \"quote\"")
}

func (cmd *AddQuoteCommand) Run(s *discordgo.Session, m *discordgo.Message) error {
	if botDatabase == nil {
		return errors.New("database not setup")
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
	botDatabase.AddQuote(db.Quote{UserID: cmd.userID, Quote: cmd.quote})
	sendMessage(fmt.Sprint("Added a new quote for ", user.Mention()), m.ChannelID, s)
	if cmd.quoteChannel != "" {
		sendMessage(fmt.Sprint("```", cmd.quote, "```", " - ", user.Mention()), cmd.quoteChannel, s)
	}
	return nil
}

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
	return fmt.Sprint(string(CommandPrefix), cmd.Name(), ` -user @user`)
}

func (cmd *GetQuoteCommand) Run(s *discordgo.Session, m *discordgo.Message) error {
	if botDatabase == nil {
		return errors.New("database not setup")
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
	quotes, err := botDatabase.GetQuotes(cmd.userID)
	if err != nil {
		return err
	}
	if len(quotes) == 0 {
		return errors.New("unable to find any quotes")
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	quoteNum := r.Int() % len(quotes)
	sendMessage(fmt.Sprint("```", quotes[quoteNum].Quote, "```\n - ", user.Mention()), m.ChannelID, s)

	return nil
}
