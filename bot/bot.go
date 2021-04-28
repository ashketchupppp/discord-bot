package bot

import (
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

type DiscordBot struct {
	Bot
	Session *discordgo.Session
	token   string
}

func (d *DiscordBot) Close() error {
	err := d.Session.Close()
	return err
}

func (d *DiscordBot) Connect() error {
	err := d.Session.Open()
	return err
}

func (d *DiscordBot) RegisterHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	d.Session.AddHandler(handler)
}

func NewDiscordBot(token string) *DiscordBot {
	ds, _ := discordgo.New("Bot " + token)
	return &DiscordBot{token: token, Session: ds}
}

func ValidUserId(userid string) bool {
	var validUserID = regexp.MustCompile(UserIDRegex)
	return validUserID.MatchString(userid)
}
