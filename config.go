package main

import (
	"io"
)

type Initialisable interface {
	// Intialises the implementor, should be called _after_ all its fields have been set
	Initialise() error
}

type Validatable interface {
	// Should check its own values to see if everything is assigned correctly
	Validate() error
}

type Configurable interface {
	// Should load data from the reader into itself
	Load(r io.Reader) error
}
