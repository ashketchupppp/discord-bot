package config

import (
	"encoding/json"
	"errors"
	"io"
)

type Validatable interface {
	// Should check its own values to see if everything is assigned correctly
	Validate() error
}

type Configuration interface {
	// Should load data from the reader into itself
	Load(r io.Reader) error
}

// A struct that stores the entire state of the application
type Bot struct {
	Token     string
	DbConnStr string
}

// Reads reader r and attempts to decode it as JSON.
// Any nonexistant keys will be ignored
func (b *Bot) Load(r io.Reader) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&b)
	if err != nil {
		return err
	}
	return nil // success!
}

// Checks the values in b to see if b is setup correctly, not everything is checked, only the important stuff
func (b *Bot) Validate() error {
	if b.Token == "" {
		return errors.New("Token not configured")
	}
	if b.DbConnStr == "" {
		return errors.New("DbConnStr not configured")
	}
	return nil
}
