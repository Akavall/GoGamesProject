package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Akavall/GoGamesProject/dice"
	"github.com/Akavall/GoGamesProject/statistics"
	"github.com/Akavall/GoGamesProject/zombie_dice"
	"github.com/nu7hatch/gouuid"
)

var templates = template.Must(template.ParseFiles("web/index.html", "web/zombie_dice.html"))
var zombie_games = make(map[*uuid.UUID]zombie_dice.GameState)

func zombie_game(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/html")

	err := templates.ExecuteTemplate(response, "zombie_dice.html", nil)

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func index(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/html")

	err := templates.ExecuteTemplate(response, "index.html", nil)

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

//TO-DO: Can start an unlimited number of games!
func start_zombie_dice(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/plain")

	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	num_players_input := request.Form["num_players"]

	var num_players int
	if len(num_players_input) == 1 {
		num_players, err = strconv.Atoi(num_players_input[0])
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Bad input on number of players! Received %d inputs, expected only 1", len(num_players_input))
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	players := make([]zombie_dice.Player, num_players)
	for i := 1; i <= num_players; i++ {
		player := "player" + strconv.Itoa(i)
		player_name_input := request.Form[player]

		var player_name string
		if len(player_name_input) == 1 {
			player_name = player_name_input[0]
		} else {
			error_message := fmt.Sprintf("Bad input on %s! Received %d inputs, expected only 1", player, len(player_name_input))
			log.Printf(error_message)
			http.Error(response, error_message, http.StatusBadRequest)
			return
		}

		is_ai_input := request.Form[player+"_ai"]

		var is_player_ai bool
		if len(is_ai_input) == 1 {
			is_player_ai, err = strconv.ParseBool(is_ai_input[0])
			if err != nil {
				http.Error(response, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			error_message := fmt.Sprintf("Bad input on AI flag for %s_ai! Received %d inputs, expected only 1", player, len(is_ai_input))
			log.Printf(error_message)
			http.Error(response, error_message, http.StatusBadRequest)
			return
		}

		score := 0
		players[i-1] = zombie_dice.Player{PlayerState: zombie_dice.InitPlayerState(), Name: player_name, IsAI: is_player_ai, TotalScore: &score}
	}

	uuid, err := uuid.NewV4()

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	game_state, err := zombie_dice.InitGameState(players)
	zombie_games[uuid] = game_state
    log.Printf("Successfully started new Zombie Dice game with ID: %s; Number of running games: %d", uuid.String(), len(zombie_games))

	fmt.Fprintf(response, "%s", uuid.String())
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

func roll_dice(response http.ResponseWriter, request *http.Request) {
	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	num_sides := 6
	sides_input := request.Form["sides"]
	if len(sides_input) == 1 {
		num_sides, _ = strconv.Atoi(sides_input[0])
	}

	my_dice := dice.InitDefaultDice(num_sides)
	side := my_dice.Roll()
	log.Printf("Rolled %d for request: \n\t%v", side.Numerical_value, request)

	fmt.Fprintf(response, "%d", side.Numerical_value)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/zombie_dice", zombie_game)
	mux.HandleFunc("/zombie_dice/start_game", start_zombie_dice)
	mux.HandleFunc("/four_dice_roll", four_dice_roll)
	mux.HandleFunc("/roll_dice", roll_dice)
	log.Printf("Started dumb Dice web server! Try it on http://localhost:8000")
	err := http.ListenAndServe(":8000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
