package db

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

type Item struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

var db *mgo.Database

func init() {
	session, err := mgo.Dial("localhost/api_db")
	if err != nil {
		fmt.Errorf("failed to connect to database: %v", err)
	}

	db = session.DB("api_db")
}

func GetAll() ([]Item, error) {
	var res []Item

	if err := db.C("items").Find(nil).All(&res); err != nil {
		return nil, fmt.Errorf("failed to fetch items: %v", err)
	}

	return res, nil
}
