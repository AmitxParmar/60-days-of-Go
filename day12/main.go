package main

import (
	"encoding/json"
	"log"
	"net/http"

	
	"github.com/cassiobotaro/60-days-of-go/day11/cards"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// future ideas:
// - json errors
// - post unique
// - persist layer
// - tests
// - separate model in diffrent file
// - model and serializers??
// controllers by package

func createCard(w http.ResponseWriter, r *http.Request) {{
	card := cards.CardSerializer{}
	err := json.NewDecoder(r.Body).Decode(&card)
	defer r.Body.Close()
	if err != nil package main
	
	func main() {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if card.Validate() {
		// not implemented yet
		card.Save()
		// set content type  as json
		// maybe in future it will turned into a middleware
	}
	}
