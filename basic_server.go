package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
	"time"

	"github.com/Akavall/GoGamesProject/dice"
	"github.com/Akavall/GoGamesProject/statistics"
	"github.com/Akavall/GoGamesProject/zombie_dice"
	"github.com/nu7hatch/gouuid"
)

const MAX_ZOMBIE_DICE_GAMES = 60

var templates = template.Must(template.ParseFiles("web/index.html", "web/zombie_dice.html", "web/zombie_dice_multi_player.html"))
var zombie_games = make(map[string]*zombie_dice.GameState)

func zombie_game(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/html")

	err := templates.ExecuteTemplate(response, "zombie_dice.html", nil)

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func zombie_game_multi_player(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/html")

	err := templates.ExecuteTemplate(response, "zombie_dice_multi_player.html", nil)

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

func start_zombie_dice(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/plain")

	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	num_players_input := request.Form["num_players"]

	log.Println("num_players_input", num_players_input)

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

	log.Println("num_players", num_players)

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
	uuid_string := uuid.String()

	if len(zombie_games) < MAX_ZOMBIE_DICE_GAMES {
		zombie_games[uuid_string] = &game_state
		log.Printf("Successfully started new Zombie Dice game with ID: %s; Number of running games: %d", uuid_string, len(zombie_games))
	} else {
		error_message := fmt.Sprintf("Maximum number of zombie dice games (%d) reached!", MAX_ZOMBIE_DICE_GAMES)
		log.Printf(error_message)
		http.Error(response, error_message, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(response, "%s", uuid_string)
}

func join_game(response http.ResponseWriter, request *http.Request) {
	log.Println("Joining the game...")
	response.Header().Set("Content-type", "text/plain")

	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	game_id_input := request.Form["game_id"]

	game_id := game_id_input[0]

	zombie_game, ok := zombie_games[game_id]

	if !ok {
		// kills the entire run
		// log.Fatalf("Cannot find a game with id: %s\n", game_id)

		log.Printf("Cannot find a game with id: %s\n", game_id)
	}

	player2_input := request.Form["player2"]
	player2_name := player2_input[0]
	score := 0
	player2 := zombie_dice.Player{PlayerState: zombie_dice.InitPlayerState(), Name: player2_name, IsAI: false, TotalScore: &score}

	(*zombie_game).Players = append((*zombie_game).Players, player2)

	log.Printf("Player: %s joined game: %s", player2_name, game_id)
}

func take_zombie_dice_turn(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/plain")

	err := request.ParseForm()
	if err != nil {
		log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	uuid, err := parse_input(request, "uuid")
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	player_name, err := parse_input(request, "player")
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("PLAYERS NAME : ", player_name)

	continue_turn_string, err := parse_input(request, "continue")
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	continue_turn, err := strconv.ParseBool(continue_turn_string)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("continue_turn_0 : %t", continue_turn)

	log.Printf("uuid : %s", uuid)
	game_state, ok := zombie_games[uuid]
 
	if !ok {
		http.Error(response, fmt.Sprintf("Game with id %s not found!", uuid), http.StatusBadRequest)
		return
	} else {
		log.Printf("Grabbed game state with id: %s", uuid)
	}

	if game_state.IsActive {
		http.Error(response, fmt.Sprintf("Game state with id %s is already active!", uuid), http.StatusBadRequest)
		return
	} else {
		game_state.IsActive = true
	}

	player_index := game_state.PlayerTurn
	active_player := game_state.Players[player_index]

	if player_name != active_player.Name {
		http.Error(response, fmt.Sprintf("%s is currently taking a turn, not %s!", active_player.Name, player_name), http.StatusBadRequest)
		game_state.IsActive = false
		return
	}

	if active_player.IsAI {
		log.Printf("Size of deck : %d\n", len(game_state.ZombieDeck.Deck.Dices))
                log.Printf("\033[034mshot: %d,brains: %d, walks: %d\033[0m", active_player.PlayerState.TimesShot, active_player.PlayerState.BrainsRolled, active_player.PlayerState.WalksTakenLastRoll)
		if zombie_dice.SimulationistAI(active_player.PlayerState.TimesShot,
			active_player.PlayerState.BrainsRolled,
		        active_player.PlayerState.WalksTakenLastRoll,
		        &game_state.ZombieDeck) == 0 {
			continue_turn = false
		}
	}

	turn_result := [3][2]string{}
	if continue_turn {
		turn_result, err = active_player.TakeTurn(&game_state.ZombieDeck)
		if err != nil {
			http.Error(response, fmt.Sprintf("Error occured while player %s was taking turn: %s", active_player.Name, err.Error()), http.StatusBadRequest)
			return
		}
	} else {
		*active_player.TotalScore += active_player.PlayerState.CurrentScore
		active_player.PlayerState.Reset()
		game_state.EndTurn()
	}

        player_turn_result := zombie_dice.PlayerTurnResult{

	TurnResult: turn_result, 
	RoundScore: active_player.PlayerState.CurrentScore,
	TimesShot: active_player.PlayerState.TimesShot,
	TotalScore: *active_player.TotalScore,
	IsDead: active_player.PlayerState.IsDead,
	Winner: game_state.Winner.Name,
	PlayerName: active_player.Name,
	ContinueTurn: continue_turn,}

	json_string, err := json.Marshal(player_turn_result)
	if err != nil {
		panic(err) //TO-DO: handle this error better
	}

	(*game_state).MoveLog = append((*game_state).MoveLog, player_turn_result)
	
	fmt.Fprintf(response, string(json_string))

	if game_state.GameOver {
		// sleeping to display the game status
		// for multi player, skiping for games with AI
		// TODO: There has to be a better way to do thi
		
		skip_sleep := false
		for _, player := range (*game_state).Players {
			if player.IsAI {
				skip_sleep = true
			}
		}
		if !skip_sleep {
			log.Println("Sleeping...")
			time.Sleep(time.Second * 30)
		}
		delete(zombie_games, uuid)
	}

	if active_player.PlayerState.IsDead {
		active_player.PlayerState.Reset()
		game_state.EndTurn()
	}

	game_state.IsActive = false
}

func get_player_turn_results(response http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		panic(err)
	}

	game_id_form, _ := request.Form["game_id"]
 	
	game_id := game_id_form[0]

	current_game, ok := zombie_games[game_id]
	if !ok {
		// This will likely to happend, if a player
		// is slow to start a game and/or another player
		// is slow to join
		log.Printf("Game id has not been found")
		return
	}

	move_log := (*current_game).MoveLog

	all_rolls := []string {}
	for _, tr := range move_log {
		roll_strings := []string {"Player: " + tr.PlayerName}
		for i := 0; i < 3; i++ {
			roll_strings = append(roll_strings, fmt.Sprintf("%s: %s", tr.TurnResult[i][0], tr.TurnResult[i][1]))
		}

		if tr.IsDead == true || tr.ContinueTurn == false {

			turn_end_string := fmt.Sprintf("Player: %s, Total Score: %d, Turn Ended\n", tr.PlayerName, tr.TotalScore)
			roll_strings = append(roll_strings, turn_end_string)
		}
		
		if tr.Winner != "" {
			winner_string := fmt.Sprintf("Winner: Player: %s", tr.Winner)

			roll_strings = append(roll_strings, winner_string) 
		}

		one_roll_string := strings.Join(roll_strings, "\n")	


		all_rolls = append(all_rolls, one_roll_string)
	}

	
	formated_moves := strings.Join(all_rolls, "\n")

	fmt.Fprintf(response, formated_moves)
	
}

func get_n_players_in_game(response http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}

	game_id_form, _ := request.Form["game_id"]
 	
	game_id := game_id_form[0]

	current_game, ok := zombie_games[game_id]
	if !ok {
		// This will likely to happend, if a player
		// is slow to start a game and/or another player
		// is slow to join
		log.Printf("Game id has not been found")
		return
	}

	fmt.Fprintf(response, "%d", len(current_game.Players))
}

func parse_input(request *http.Request, field string) (s string, err error) {
	
	input_array := request.Form[field]
	parsed_input := ""
	if len(input_array) == 1 {
		parsed_input = input_array[0]
	} else {
		error_message := fmt.Sprintf("Bad input on %s! Received %d inputs, expected only 1", field, len(input_array))
		log.Printf(error_message)
		return "", errors.New(error_message)
	}
	return parsed_input, nil
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
	mux.HandleFunc("/zombie_dice_multi_player", zombie_game_multi_player)
	mux.HandleFunc("/zombie_dice/start_game", start_zombie_dice)
	mux.HandleFunc("/zombie_dice_multi_player/start_game", start_zombie_dice)

	mux.HandleFunc("/zombie_dice_multi_player/join_game", join_game)
	mux.HandleFunc("/zombie_dice/take_turn", take_zombie_dice_turn)
	mux.HandleFunc("/zombie_dice_multi_player/take_turn", take_zombie_dice_turn)
	mux.HandleFunc("/zombie_dice_multi_player/get_player_turn_results", get_player_turn_results)
	mux.HandleFunc("/zombie_dice_multi_player/get_n_players_in_game", get_n_players_in_game)

	mux.HandleFunc("/four_dice_roll", four_dice_roll)
	mux.HandleFunc("/roll_dice", roll_dice)

	log.Printf("Started dumb Dice web server! Try it on http://ip_address:8000")

	err := http.ListenAndServe("0.0.0.0:8000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
