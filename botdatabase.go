package main

type BotDatabase interface {
	GetQuotes(userID string) ([]Quote, error)
	AddQuote(q Quote) error
	Connect() error
}

type Quote struct {
	UserID string
	Quote  string
}
