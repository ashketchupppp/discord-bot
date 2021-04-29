package main

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var (
	UserIDRegex = `<@!(\d{18})>`
)

type Bot interface {
	Connect() error
	Close() error
	RegisterListener(s *discordgo.Session, m *discordgo.MessageCreate) error
	SendMessage(c string, msg string) error
}

// Stores the state of the entire application
type DiscordBot struct {
	Session   *discordgo.Session
	Token     string
	DbConnStr string
}

// Closes the discord session
func (d *DiscordBot) Close() error {
	err := d.Session.Close()
	return err
}

// Opens the discord session
func (d *DiscordBot) Connect() error {
	err := d.Session.Open()
	return err
}

// Registers a function handler with the discord session
func (d *DiscordBot) RegisterHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	d.Session.AddHandler(handler)
}

func NewDiscordBot(token string) *DiscordBot {
	ds, _ := discordgo.New("Bot " + token)
	return &DiscordBot{Token: token, Session: ds}
}

func ValidUserId(userid string) bool {
	var validUserID = regexp.MustCompile(UserIDRegex)
	return validUserID.MatchString(userid)
}

// Reads reader r and attempts to decode it as JSON.
// Any JSON keys that aren't struct fields (and vice versa) will be ignored
func (b *DiscordBot) Load(r io.Reader) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&b)
	if err != nil {
		return err
	}
	return nil // success!
}

// Checks the values in b to see if b is setup correctly, not everything is checked, only the important stuff
func (b *DiscordBot) Validate() error {
	if b.Token == "" {
		return errors.New("Token not configured")
	}
	if b.DbConnStr == "" {
		return errors.New("DbConnStr not configured")
	}
	return nil
}
