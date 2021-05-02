package main

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	*mongo.Client
	ConnStr         string
	QuoteCollection string
	DBName          string
}

func (m *MongoDB) Validate() error {
	if m.ConnStr == "" {
		return errors.New("DbConnStr not configured")
	}
	return nil
}

func (m *MongoDB) Connect() error {
	copt := options.Client().ApplyURI(m.ConnStr)
	client, err := mongo.Connect(context.Background(), copt)
	if err != nil {
		return err
	}
	if m.QuoteCollection == "" {
		m.QuoteCollection = "quotes"
	}
	if m.DBName == "" {
		m.DBName = "discordbot"
	}
	m.Client = client
	return nil
}

func (m *MongoDB) AddQuote(q Quote) error {
	collection := m.Database(m.DBName).Collection(m.QuoteCollection)
	_, err := collection.InsertOne(context.TODO(), q)
	return err
}

func (m *MongoDB) GetQuotes(userID string) ([]Quote, error) {
	collection := m.Database(m.DBName).Collection(m.QuoteCollection)
	filter := bson.D{{"user", userID}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	quotes := []Quote{}
	err = cursor.All(context.TODO(), &quotes)
	if err != nil {
		return nil, err
	}
	return quotes, err
}
