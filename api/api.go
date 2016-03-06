package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mycodesmells/mongo-go-api/db"
)

func handleError(err error, message string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(message, err)))
}

func handleCustomError(message string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

// GetAllItems returns a list of all database items to the response.
func GetAllItems(w http.ResponseWriter, req *http.Request) {
	rs, err := db.GetAll()
	if err != nil {
		handleError(err, "Failed to load database items: %v", w)
		return
	}

	bs, err := json.Marshal(rs)
	if err != nil {
		handleError(err, "Failed to load marshal data: %v", w)
		return
	}

	w.Write(bs)
}

// GetItem returns a single database item matching given ID parameter.
func GetItem(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	rs, err := db.GetOne(id)
	if err != nil {
		handleError(err, "Failed to read database: %v", w)
		return
	}

	bs, err := json.Marshal(rs)
	if err != nil {
		handleError(err, "Failed to marshal data: %v", w)
		return
	}

	w.Write(bs)
}

// PostItem saves an item (form data) into the database.
func PostItem(w http.ResponseWriter, req *http.Request) {
	ID := req.FormValue("id")
	valueStr := req.FormValue("value")
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		handleError(err, "Failed to parse input data: %v", w)
		return
	}

	item := db.Item{ID: ID, Value: value}

	if err = db.Save(item); err != nil {
		handleError(err, "Failed to save data: %v", w)
		return
	}

	w.Write([]byte("OK"))
}

// DeleteItem removes a single item (identified by parameter) from the database.
func DeleteItem(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	if err := db.Remove(id); err != nil {
		handleError(err, "Failed to remove item: %v", w)
		return
	}

	w.Write([]byte("OK"))
}

// GetImage returns an image bytes.
func GetImage(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	itemID := vars["id"]
	imageID := vars["img"]

	image, err := db.LoadImage(itemID, imageID)
	if err != nil {
		handleError(err, "Failed to load image: %v", w)
	}

	w.Write(image)
}

// UploadImage saves an image.
func UploadImage(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	itemID := vars["id"]

	file, headers, err := req.FormFile("image")
	if err != nil {
		handleError(err, "Failed to process image upload: %v", w)
		return
	}

	if file == nil {
		handleCustomError("File parameter is missing", w)
		return
	}
	defer file.Close()

	if err = db.SaveImage(itemID, headers, file); err != nil {
		handleError(err, "Failed to save image: %v", w)
		return
	}

	w.Write([]byte("OK"))
}

// DeleteImage removes an image.
func DeleteImage(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	itemID := vars["id"]
	imageID := vars["img"]

	if err := db.RemoveImage(itemID, imageID); err != nil {
		handleError(err, "Failed to delete image: %v", w)
		return
	}

	w.Write([]byte("OK"))
}
