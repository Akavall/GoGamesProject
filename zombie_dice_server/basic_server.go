package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Akavall/GoGamesProject/dynamo_db_tools"
	"github.com/Akavall/GoGamesProject/zombie_dice"
	"github.com/nu7hatch/gouuid"
)

var templates = template.Must(template.ParseFiles("zombie_dice.html"))

var zombie_games = make(map[string]*zombie_dice.GameState)

type id_to_name_type map[string]string

var id_to_name id_to_name_type = map[string]string{}

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
		player_id_input := request.Form[player]

		var player_id string
		if len(player_id_input) == 1 {
			player_id = player_id_input[0]
		} else {
			error_message := fmt.Sprintf("Bad input on %s! Received %d inputs, expected only 1", player, len(player_id_input))
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
		players[i-1] = zombie_dice.Player{PlayerState: zombie_dice.InitPlayerState(), Id: player_id, IsAI: is_player_ai, TotalScore: &score}
	}

	uuid, err := uuid.NewV4()

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	uuid_string := uuid.String()
	game_state, err := zombie_dice.InitGameState(players, uuid_string)

	err = dynamo_db_tools.PutGameStateInDynamoDB(game_state)

	if err != nil {
		log.Println("Was not able to put GameState in DynamoDB", err)
	} else {
		log.Printf("Put GateState associated with %s in DynamoDB table: GameStates", uuid_string)
	}

	fmt.Fprintf(response, "%s", uuid_string)
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

	player_id, err := parse_input(request, "player")
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("PLAYERS NAME : ", player_id)

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

	game_state, err := dynamo_db_tools.GetGameStateFromDynamoDB(uuid)

	if err != nil {
		http.Error(response, fmt.Sprintf("Game with id %s not found!, %v", uuid, err), http.StatusBadRequest)
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

	if player_id != active_player.Id {
		http.Error(response, fmt.Sprintf("%s is currently taking a turn, not %s!", active_player.Id, player_id), http.StatusBadRequest)
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

	turn_result := make([][]string, 3)

	for i := 0; i < 3; i++ {
		turn_result[i] = make([]string, 2)
	}

	if continue_turn {
		turn_result, err = active_player.TakeTurn(&game_state.ZombieDeck)
		if err != nil {
			http.Error(response, fmt.Sprintf("Error occured while player %s was taking turn: %s", active_player.Id, err.Error()), http.StatusBadRequest)
			return
		}
	} else {
		*active_player.TotalScore += active_player.PlayerState.CurrentScore
		active_player.PlayerState.Reset()
		game_state.EndTurn()
	}

	if active_player.PlayerState.IsDead {
		active_player.PlayerState.Reset()
		game_state.EndTurn()
	}

	log.Printf("Winner: %s", game_state.Winner.Id)

	player_turn_result := zombie_dice.PlayerTurnResult{

		TurnResult:   turn_result,
		RoundScore:   active_player.PlayerState.CurrentScore,
		TimesShot:    active_player.PlayerState.TimesShot,
		TotalScore:   *active_player.TotalScore,
		IsDead:       active_player.PlayerState.IsDead,
		Winner:       game_state.Winner.Id,
		PlayerId:     active_player.Id,
		ContinueTurn: continue_turn}

	json_string, err := json.Marshal(player_turn_result)
	if err != nil {
		panic(err) //TO-DO: handle this error better
	}

	game_state.MoveLog = append(game_state.MoveLog, player_turn_result)

	fmt.Fprintf(response, string(json_string))

	log.Printf("game_state.GameOver: %t", game_state.GameOver)

	if game_state.GameOver {
		// sleeping to display the game status
		// for multi player, skiping for games with AI
		// TODO: There has to be a better way to do this

		skip_sleep := false
		for _, player := range game_state.Players {
			if player.IsAI {
				skip_sleep = true
			}
		}
		if !skip_sleep {
			log.Println("Sleeping...")
			time.Sleep(time.Second * 30)
		}

		err := dynamo_db_tools.DeleteGameStateFromDynamoDB(uuid)

		if err != nil {
			log.Println("Was not able to delete GameState in DynamoDB, uuid: %s", err, uuid)
		} else {
			log.Printf("Deleted GameState associated with %s in DynamoDB table: GameStates", uuid)
		}
	}

	game_state.IsActive = false

	if !game_state.GameOver {

		err = dynamo_db_tools.PutGameStateInDynamoDB(game_state)

		if err != nil {
			log.Println("Was not able to put GameState in DynamoDB", err)
		} else {
			log.Printf("Put/update GameState associated with %s in DynamoDB table: GameStates", uuid)
		}
	}
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

func main() {

	f, err := os.OpenFile("/var/log/ZombieDice/logfile.txt", os.O_RDWR|os.O_APPEND, 0660)

	if err != nil {
		fmt.Println("Could not open logfile.txt")
	}

	log.SetOutput(f)

	mux := http.NewServeMux()
	mux.HandleFunc("/zombie_dice", zombie_game)
	mux.HandleFunc("/zombie_dice/start_game", start_zombie_dice)

	mux.HandleFunc("/zombie_dice/take_turn", take_zombie_dice_turn)

	log.Printf("Started dumb Dice web server! Try it on http://ip_address:8000")

	err = http.ListenAndServe("0.0.0.0:8000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
