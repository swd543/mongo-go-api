package api

import (
	"net/http"

	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mycodesmells/mongo-go-api/db"
)

func Items(w http.ResponseWriter, req *http.Request) {
	rs, err := db.GetAll()
	if err != nil {
		fmt.Errorf("AAAA %v", err)
	}

	bs, err := json.Marshal(rs)
	if err != nil {
		fmt.Errorf("BBBB %v", err)
	}

	w.Write(bs)
}

func Item(w http.ResponseWriter, req *http.Request) {
	rs, err := db.GetAll()
	if err != nil {
		fmt.Errorf("AAAA %v", err)
	}

	_, err = json.Marshal(rs)
	if err != nil {
		fmt.Errorf("BBBB %v", err)
	}
	vars := mux.Vars(req)
	id := vars["id"]

	w.Write([]byte(id))
}
