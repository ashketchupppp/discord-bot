package config

import (
	"strings"
	"testing"
)

// Tests that the Load method is able to load valid JSON data from a reader into a struct
func TestLoadConfigValidData(t *testing.T) {
	reader := strings.NewReader(`{"Token" : "atoken", "DbConnStr" : "dbconnstr"}`)
	b := &Bot{}
	err := b.Load(reader)
	if err != nil {
		t.Errorf(err.Error())
	}
	// Read a few values to make sure things are read correctly
	// Don't need to read everything, if some can then they all can
	if b.Token != "atoken" {
		t.Errorf("b.Token != \"atoken\"")
	}
	if b.DbConnStr != "dbconnstr" {
		t.Errorf("b.DbConnStr != \"dbconnstr\"")
	}
}

// Tests that the Load method is able to load valid JSON data from a reader into a struct
func TestLoadConfigMissingDataCausesInvalid(t *testing.T) {
	reader := strings.NewReader(`{"DbConnStr" : "dbconnstr"}`)
	b := &Bot{}
	err := b.Load(reader)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = b.Validate()
	if err == nil {
		t.Errorf("b.Validate threw even when Token was missing")
	}
}
