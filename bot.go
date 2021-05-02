package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mattn/go-shellwords"
)

var (
	UserIDRegex = `<@!(\d{18})>`
)

var (
	availableCommands map[string]Command
	availableFeatures []string
)

func init() {
	availableCommands = make(map[string]Command)

	availableCommands[HelpCmdName] = &HelpCommand{name: HelpCmdName}
	availableCommands[AddQuoteCmdName] = &AddQuoteCommand{name: AddQuoteCmdName}
	availableCommands[GetQuoteCmdName] = &GetQuoteCommand{name: GetQuoteCmdName}

	availableFeatures = []string{
		LeaveMessageFeatureName,
	}

}

type Bot interface {
	Connect() error
	Disconnect() error
	SendMessage(c string, msg string) error
	Database() *BotDatabase
	ParseCommand() (*Command, error)
	GetCommand() (*Command, error)
}

// Stores the state of the entire application
type DiscordBot struct {
	Session       *discordgo.Session
	Token         string
	MongoDatabase *MongoDB
	CommandPrefix string

	Settings        map[string]string
	EnabledCommands []string
	EnabledFeatures []string
}

// Closes the discord session
func (d *DiscordBot) Disconnect() error {
	err := d.Session.Close()
	return err
}

// Opens the discord session
func (b *DiscordBot) Connect() error {
	err := b.Session.Open()
	return err
}

func (b *DiscordBot) SendError(e error, chanID string, s *discordgo.Session) {
	b.SendMessage(fmt.Sprint("```", e.Error(), "\nfor help, use $help", "```"), chanID, s)
}

func (b *DiscordBot) SendMessage(msg, chanID string, s *discordgo.Session) {
	s.ChannelMessageSend(chanID, msg)
}

func ValidUserId(userid string) bool {
	var validUserID = regexp.MustCompile(UserIDRegex)
	return validUserID.MatchString(userid)
}

// Returns a database to use
func (b *DiscordBot) Database() (BotDatabase, error) {
	if b.MongoDatabase != nil {
		return b.MongoDatabase, nil
	}
	return nil, errors.New("unable to determine database to use")
}

func (b *DiscordBot) GetSetting(name string) (string, error) {
	if val, exists := b.Settings[name]; exists {
		return val, nil
	} else {
		return "", fmt.Errorf(name, " has not been set")
	}
}

// Reads reader r and attempts to decode it as JSON.
// Any JSON keys that aren't struct fields (and vice versa) will be ignored
// Also creates the command objects and puts them in the availableCommands map
func (b *DiscordBot) Load(r io.Reader) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&b)
	if err != nil {
		return err
	}
	return nil // success!
}

// Checks the values in b to see if b is setup correctly,
// all Validatable objects the DiscordBot stores will also be validated
func (b *DiscordBot) Validate() error {
	var err error
	if b.Token == "" {
		return errors.New("token not configured")
	}

	// Check that all enabled commands are actually commands
	for i := range b.EnabledCommands {
		if _, cmdIsAvailable := availableCommands[b.EnabledCommands[i]]; !cmdIsAvailable {
			return fmt.Errorf(b.EnabledCommands[i], " is not a command")
		}
	}

	// get the commands to validate anything they need to
	for i := range b.EnabledCommands {
		if cmd, cmdIsAvailable := availableCommands[b.EnabledCommands[i]]; cmdIsAvailable {
			err = cmd.Validate()
			if err != nil {
				return err
			}
		}
	}

	// check that all enabled features actually exist
	for i := range b.EnabledFeatures {
		featureExists := false
		for _, y := range availableFeatures {
			if y == b.EnabledFeatures[i] {
				featureExists = true
			}
		}
		if !featureExists {
			return fmt.Errorf(b.EnabledFeatures[i], " is not a feature, it must be one of ", availableFeatures)
		}
	}

	// Validate the database is setup correctly
	if b.MongoDatabase != nil {
		err = b.MongoDatabase.Validate()
	}
	if err != nil {
		return err
	}

	// We should have a database setup by now
	var db BotDatabase
	db, err = b.Database()
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("unable to determine bot database")
	}
	return nil
}

func (b *DiscordBot) CommandIsEnabled(name string) bool {
	cmdEnabled := false
	for _, v := range b.EnabledCommands {
		if v == name {
			cmdEnabled = true
		}
	}
	return cmdEnabled
}

func (b *DiscordBot) GetCommand(cmdName string) (Command, error) {
	cmd, cmdExists := availableCommands[cmdName]
	if !cmdExists {
		return nil, errors.New("unable to find command")
	}
	if !b.CommandIsEnabled(cmd.Name()) {
		return nil, fmt.Errorf(cmd.Name(), " is not enabled")
	}
	return cmd, nil
}

func (b *DiscordBot) ParseCommand(cmdStr string) (Command, error) {
	// if the command contains user ids (the "@Username" string), they need quotes around them
	re := regexp.MustCompile(UserIDRegex)
	userIDs := re.FindAll([]byte(cmdStr), 10)
	for id := range userIDs {
		cmdStr = strings.Replace(cmdStr, string(userIDs[id]), fmt.Sprint("\"", string(userIDs[id]), "\""), -1)
	}

	args, err := shellwords.Parse(cmdStr)
	if err != nil {
		return nil, err
	}
	cmd, err := b.GetCommand(args[0])
	if err != nil {
		return nil, err
	}
	fs := cmd.FlagSet()
	err = fs.Parse(args[1:])
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (b *DiscordBot) Initialise() error {
	// Setup the discord session
	var err error
	b.Session, err = discordgo.New("Bot " + b.Token)
	if err != nil {
		return err
	}

	// Setup the event handlers
	b.Session.AddHandler(NewMessageHandler)
	b.Session.AddHandler(LeaveMessageHandler)

	for f := range b.EnabledFeatures {
		if b.EnabledFeatures[f] == LeaveMessageFeatureName {
			// b.Session.AddHandler(LeaveMessageHandler)
		}
	}

	// Setup our database
	db, err := b.Database()
	if err != nil {
		return err
	}
	err = db.Connect()
	if err != nil {
		return err
	}
	return nil
}

func LeaveMessageHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	bot := GetDiscordBot()
	leaveChannel, err := bot.GetSetting(LeaveChannelSettingName)
	if err != nil {
		log.Print(err.Error())
	}
	bot.SendMessage(fmt.Sprintf(m.Mention(), " left the server!"), leaveChannel, s)
}

func NewMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	bot := GetDiscordBot()
	if m.Author.ID == s.State.User.ID || m.Content[0] != []byte(bot.CommandPrefix)[0] {
		return
	}
	m.Content = m.Content[1:]
	cmd, err := bot.ParseCommand(m.Content)
	if err != nil {
		bot.SendError(err, m.ChannelID, s)
		log.Print(err.Error())
		return
	}
	err = cmd.Run(s, m.Message)
	if err != nil {
		bot.SendError(err, m.ChannelID, s)
		log.Print(err.Error())
	}
}
