package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mycodesmells/mongo-go-api/api"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/items", api.Items).Methods("GET")
	r.HandleFunc("/api/items/{id}", api.Item).Methods("GET")

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":3000", r)
}
