package db

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Item represents a sample database entity.
type Item struct {
	ID     string   `json:"id" bson:"_id,omitempty"`
	Value  int      `json:"value"`
	Images []string `json:"images"`
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
	res := Item{}
	if err := collection().Find(bson.M{"_id": id}).One(&res); err != nil {
		return err
	}
	for _, imageID := range res.Images {
		if err := RemoveImage(id, imageID); err != nil {
			fmt.Printf("Failed to remove image %s_%s: %v", id, imageID, err)
		}
	}

	return collection().Remove(bson.M{"_id": id})
}

// LoadImage returns an image from the disk.
func LoadImage(itemID, imageID string) ([]byte, error) {
	filename := fmt.Sprintf("%s_%s", itemID, imageID)
	return ioutil.ReadFile("./media/upload/" + filename)
}

// SaveImage saves an image to the disk and updates an Item in the database.
func SaveImage(itemID string, headers *multipart.FileHeader, file multipart.File) error {
	filename := fmt.Sprintf("%s_%s", itemID, headers.Filename)
	saved, err := os.Create("./media/upload/" + filename)
	if err != nil {
		return err
	}
	defer saved.Close()

	_, err = io.Copy(saved, file)
	if err != nil {
		return err
	}

	query := bson.M{"_id": itemID}
	res := Item{}
	if err := collection().Find(query).One(&res); err != nil {
		return err
	}
	res.Images = append(res.Images, headers.Filename)
	return collection().Update(query, res)
}

// RemoveImage deletes an image from the disk and updates an Item in the database.
func RemoveImage(itemID, imageID string) error {
	filename := fmt.Sprintf("%s_%s", itemID, imageID)
	if err := os.Remove("./media/upload/" + filename); err != nil {
		return err
	}

	query := bson.M{"_id": itemID}
	res := Item{}
	if err := collection().Find(query).One(&res); err != nil {
		return err
	}

	removeImageFromItem(&res, imageID)

	return collection().Update(query, res)
}

func removeImageFromItem(it *Item, imageID string) {
	for i, name := range it.Images {
		if name == imageID {
			it.Images = it.Images[:i+copy(it.Images[i:], it.Images[i+1:])]
			return
		}
	}
}
