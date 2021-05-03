package main

import (
	"context"
	"errors"
	"log"

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
	log.Println("attempting to connect to ", m.ConnStr)
	copt := options.Client().ApplyURI(m.ConnStr)
	client, err := mongo.Connect(context.Background(), copt)
	if err != nil {
		log.Println("connection failed: ", err.Error())
		return err
	}
	log.Println("connection successful")
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
	log.Println("adding quote ", q)
	collection := m.Database(m.DBName).Collection(m.QuoteCollection)
	_, err := collection.InsertOne(context.TODO(), q)
	if err == nil {
		log.Println("failed to add quote: ", err.Error())
	} else {
		log.Println("successfully added quote", q)
	}
	return err
}

func (m *MongoDB) GetQuotes(userID string) ([]Quote, error) {
	log.Println("getting quotes for ", userID)
	collection := m.Database(m.DBName).Collection(m.QuoteCollection)
	filter := bson.D{{"user", userID}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("error occurred when getting quotes for ", userID, " : ", err.Error())
		return nil, err
	}
	quotes := []Quote{}
	err = cursor.All(context.TODO(), &quotes)
	if err != nil {
		return nil, err
	}
	return quotes, err
}
