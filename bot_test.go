package main

import (
	"strings"
	"testing"
)

/*
Why doesnt this work?
func TestCanAddHandler(t *testing.T) {
	b := NewDiscordBot("a token")
	b.Session.AddHandler(RunTestCommand)
}
*/

// Tests that the Load method is able to load valid JSON data from a reader into a struct
func TestBotLoadValidData(t *testing.T) {
	reader := strings.NewReader(`{"Token" : "atoken", "mongoDatabase" : {"connstr" : "test"}}`)
	b := &DiscordBot{}
	err := b.Load(reader)
	if err != nil {
		t.Errorf(err.Error())
	}
	// Read a few values to make sure things are read correctly
	// Don't need to read everything, if some can then they all can
	if b.Token != "atoken" {
		t.Errorf("b.Token != \"atoken\"")
	}
	if b.MongoDatabase == nil {
		t.Errorf("b.MongoDatabase == nil")
	}
	if b.MongoDatabase.ConnStr != "test" {
		t.Errorf("b.MongoDatabase.ConnStr != test")
	}
}

// Tests that the Load method is able to load valid JSON data from a reader into a struct
func TestBotLoadMissingDataCausesInvalid(t *testing.T) {
	reader := strings.NewReader(`{}`)
	b := &DiscordBot{}
	err := b.Load(reader)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = b.Validate()
	if err == nil {
		t.Errorf("b.Validate threw even when Token was missing")
	}
}
