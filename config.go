package main

import (
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
