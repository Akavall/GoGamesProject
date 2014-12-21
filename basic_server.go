package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-temlyakov/GoGamesProject/dice"
	"github.com/a-temlyakov/GoGamesProject/statistics"
)

func dice_roll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/plain")

	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	my_dice := dice.InitDefaultDice(6)
	side := my_dice.Roll()
	log.Printf("Rolled: %d for request...", side.Numerical_value)

	// Actual response sent to web client
	fmt.Fprintf(response, "%d", side.Numerical_value)
}

func four_dice_roll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/plain")

	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	score := dice.N_dice_roll(6, 4)
	roll_prob, prob_lower, prob_higher := statistics.CalcRollProbabilities(score, 4, 6)

	log.Printf("Rolled: %d for request...", score)

	// Actual response sent to web client
	fmt.Fprintf(response, " score : %d\n roll prob : %f\n prob lower : %f\n prob higher : %f\n", score, roll_prob, prob_lower, prob_higher)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", four_dice_roll)
	log.Printf("Started dumb Dice web server! Try it on http://localhost:8000")
	err := http.ListenAndServe(":8000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
