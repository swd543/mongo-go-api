package db

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Item represents a sample database entity.
type Item struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Value int    `json:"value"`
}

var db *mgo.Database

func init() {
	session, err := mgo.Dial("localhost/api_db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db = session.DB("api_db")
}

func collection() *mgo.Collection {
	return db.C("items")
}

// GetAll returns all items from the database.
func GetAll() ([]Item, error) {
	res := []Item{}

	if err := collection().Find(nil).All(&res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetOne returns a single item from the database.
func GetOne(id string) (*Item, error) {
	res := Item{}

	if err := collection().Find(bson.M{"_id": id}).One(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Save inserts an item to the database.
func Save(item Item) error {
	return collection().Insert(item)
}

// Remove deletes an item from the database
func Remove(id string) error {
	return collection().Remove(bson.M{"_id": id})
}
