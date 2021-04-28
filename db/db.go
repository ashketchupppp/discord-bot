package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuoteProvider interface {
	GetQuotes(userID string) (string, error)
}

type QuoteStorer interface {
	AddQuote(userID, quote string) error
}

type BotDB struct {
	*mongo.Client
	quoteCollection string
	dbname          string
}

type Quote struct {
	UserID string
	Quote  string
}

func (m *BotDB) AddQuote(q Quote) error {
	collection := m.Database(m.dbname).Collection(m.quoteCollection)
	_, err := collection.InsertOne(context.TODO(), q)
	return err
}

func (m *BotDB) GetQuotes(userID string) ([]Quote, error) {
	collection := m.Database(m.dbname).Collection(m.quoteCollection)
	cursor, err := collection.Find(context.TODO(), bson.D{{"UserID", userID}})
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

func (m *BotDB) SetQuoteCollection(coll string) {
	m.quoteCollection = coll
}

func (m *BotDB) SetDatabaseName(name string) {
	m.dbname = name
}

func NewBotDB(connStr string) (*BotDB, error) {
	copt := options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(context.Background(), copt)
	if err != nil {
		return nil, err
	}
	mdb := &BotDB{}
	mdb.SetDatabaseName("discordbots")
	mdb.SetQuoteCollection("quotes")
	mdb.Client = client
	return mdb, nil
}
