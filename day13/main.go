package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	valid "github.com/asaskevich/govalidator"
	"github.com/cassiobotaro/60-days-of-go/day13/cards"
	"github.com/cassiobotaro/60-days-of-go/day13/database"
	"github.com/gorilla/mux"
)

// future ideas:
// - paginate results
// - tests
// controllers by package

// Ugly but for while is the solution
var db = database.NewMemoryDB()

// RenderJSON render a content as json(thinking about middleware)
func RenderJSON(w http.ResponseWriter, content interface{}, statusCode int) {
	// Set Content-Type as json
	w.Header().Set("Content-Type", "application/json; charset-utf-8")
	// HTTP STATUS CODE
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(content)
	if err != nil {
		log.Println(err)
	}
}

func createCard(w http.ResponseWriter, r *http.Request) {
	// initialize a card
	card := cards.Card{}
	// decode received content into struct
	err := json.NewDecoder(r.Body).Decode(&card)
	defer r.Body.Close()
	if err != nil {
		// STATUS 422 - Unprocessable entity
		renderJSON(w, map[string]string{"errors": err.Error()}, http.StatusUnprocessableEntity)
		return
	}
	//if is a valid card
	result, err := valid.ValidateStruct(card)
	if result {
		// create card
		db.CreateCard(&card)
		RenderJSON(w, card, http.StatusCreated)
	} else {
		// STATUS 401 - BAD REQUEST
		RenderJSON(w, map[string]string{"errors": err.Error()}, http.StatusBadRequest)
	}
}

func allCards(w http.ResponseWriter, r *http.Request) {
	//list all cards
	cardList := db.AllCards()
	RenderJSON(w, cardList, http.StatusOK)
}

func getCard(w http.ResponseWriter, r *http.Request) {
	// Get the id from path
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		RenderJSON(w, err, http.StatusInternalServerError)
		return
	}

	//get the card by id
	card, err := db.getCard(id)
	switch err {
	case database.ErrCardNotFound:
		RenderJSON(w, err, http.StatusNotFound)
	case nil:
		RenderJSON(w, card, http.StatusOK)
	default:
		RenderJSON(w, err, http.StatusInternalServerError)
	}
}
