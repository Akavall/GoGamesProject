package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Akavall/GoGamesProject/dice"
	"github.com/Akavall/GoGamesProject/statistics"
)

func dice_roll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/plain")

	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	num_sides := 6
	sides_input := request.Form["sides"]
	if len(sides_input) == 1 {
		num_sides, _ = strconv.Atoi(sides_input[0])
	}

	my_dice := dice.InitDefaultDice(num_sides)
	side := my_dice.Roll()
	log.Printf("Rolled %d for request: \n\t%v", side.Numerical_value, request)

	// Actual response sent to web client
	fmt.Fprintf(response, "\nRolling dice with %d sides:\n", num_sides)
	fmt.Fprintf(response, "%d", side.Numerical_value)
}

func four_dice_roll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/plain")

	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	num_sides := 6
	sides_input := request.Form["sides"]
	if len(sides_input) == 1 {
		num_sides, _ = strconv.Atoi(sides_input[0])
	}

	roll_times := 4
	roll_times_input := request.Form["rolltimes"]
	if len(roll_times_input) == 1 {
		roll_times, _ = strconv.Atoi(roll_times_input[0])
	}

	score := dice.InitDefaultDice(num_sides).RollNTimes(roll_times).SumSides()

	roll_prob, prob_lower, prob_higher := statistics.CalcRollProbabilities(score, roll_times, num_sides)

	log.Printf("Rolled %d for request: \n\t%v", score, request)

	// Actual response sent to web client
	fmt.Fprintf(response, "\nRolling dice with %d sides %d times:\n", num_sides, roll_times)
	fmt.Fprintf(response, " score : %d\n roll prob : %f\n prob lower : %f\n prob higher : %f\n", score, roll_prob, prob_lower, prob_higher)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/dice_roll", dice_roll)
	mux.HandleFunc("/four_dice_roll", four_dice_roll)
	log.Printf("Started dumb Dice web server! Try it on http://localhost:8000")
	err := http.ListenAndServe(":8000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
