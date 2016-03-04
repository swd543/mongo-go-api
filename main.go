package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mycodesmells/mongo-go-api/api"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/items", api.GetAllItems).Methods("GET")
	r.HandleFunc("/api/items/{id}", api.GetItem).Methods("GET")
	r.HandleFunc("/api/items", api.PostItem).Methods("POST")
	r.HandleFunc("/api/items/{id}", api.DeleteItem).Methods("DELETE")

	http.ListenAndServe(":3000", r)
}
