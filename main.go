package main

import (
	"coffeeserver/models"
	storage "coffeeserver/storagev2"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var store models.Storer

// go get -u github.com/gorilla/mux
func main() {
	s := storage.New()
	// set up the validator to check if the values in the map are a reference to a coffee or not
	s.ValidateFn = func(v interface{}) bool {
		_, ok := v.(models.Coffee)
		return ok
	}

	defer s.Close()
	store = s

	router := mux.NewRouter()
	router.HandleFunc("/coffee", handleGetCoffees).Methods(http.MethodGet)
	router.HandleFunc("/coffee", handleMakeCoffee).Methods(http.MethodPost)
	log.Fatalln(http.ListenAndServe(":3000", router))
}

func handleGetCoffees(w http.ResponseWriter, r *http.Request) {
	coffees := store.State()

	bs, err := json.Marshal(coffees)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func handleMakeCoffee(w http.ResponseWriter, r *http.Request) {
	// check the body

	var coffee models.Coffee

	if err := json.NewDecoder(r.Body).Decode(&coffee); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if len(coffee.Name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := store.Val(coffee.Name); err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err := store.Add(coffee.Name, coffee); err != nil {
		log.Printf("error adding: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
