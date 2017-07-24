package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Akavall/GoGamesProject/dice"
	"github.com/Akavall/GoGamesProject/dynamo_db_tools"
	"github.com/Akavall/GoGamesProject/statistics"
	"github.com/Akavall/GoGamesProject/zombie_dice"
	"github.com/nu7hatch/gouuid"
)

const MAX_ZOMBIE_DICE_GAMES = 60

var templates = template.Must(template.ParseFiles("web/index.html", "web/zombie_dice.html", "web/zombie_dice_multi_player.html"))
var zombie_games = make(map[string]*zombie_dice.GameState)
var zombie_chats = make(map[string]*zombie_dice.ZombieChat)

type id_to_name_type map[string]string

var id_to_name id_to_name_type = map[string]string{}

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

	if len(zombie_games) < MAX_ZOMBIE_DICE_GAMES {
		zombie_games[uuid_string] = &game_state
		log.Printf("Successfully started new Zombie Dice game with ID: %s; Number of running games: %d", uuid_string, len(zombie_games))
	} else {
		error_message := fmt.Sprintf("Maximum number of zombie dice games (%d) reached!", MAX_ZOMBIE_DICE_GAMES)
		log.Printf(error_message)
		http.Error(response, error_message, http.StatusBadRequest)
		return
	}

	zombie_chat := zombie_dice.ZombieChat{}

	if len(zombie_chats) < MAX_ZOMBIE_DICE_GAMES {
		zombie_chats[uuid_string] = &zombie_chat
		log.Printf("Successfully started new Zombie Dice chat with ID: %s; Number of running chats: %d", uuid_string, len(zombie_chats))
	} else {
		error_message := fmt.Sprintf("Maximum number of zombie dice chats (%d) reached!", MAX_ZOMBIE_DICE_GAMES)
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
	player2 := zombie_dice.Player{PlayerState: zombie_dice.InitPlayerState(), Id: player2_name, IsAI: false, TotalScore: &score}

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
	// Step I
	// We would need to take the game_state from DynamoDB
	// game_state, ok := zombie_games[uuid]
	fmt.Printf("uuid: %s\n", uuid)
	game_state, err := dynamo_db_tools.GetGameStateFromDynamoDB(uuid)
	fmt.Printf("game_state: %v\n", game_state)

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

	turn_result := [3][2]string{}
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

		// Step III delete game_state here
		// Need to delete game from dynamoDB table
		// delete(zombie_games, uuid)

		err := dynamo_db_tools.DeleteGameStateFromDynamoDB(uuid)

		if err != nil {
			log.Println("Was not able to delete GameState in DynamoDB, uuid: %s", err, uuid)
		} else {
			log.Printf("Deleted GameState associated with %s in DynamoDB table: GameStates", uuid)
		}

		delete(zombie_chats, uuid)
	}

	if active_player.PlayerState.IsDead {
		active_player.PlayerState.Reset()
		game_state.EndTurn()
	}

	game_state.IsActive = false

	//Step II, we can save the game_state here, just another put
	err = dynamo_db_tools.PutGameStateInDynamoDB(game_state)

	if err != nil {
		log.Println("Was not able to put GameState in DynamoDB", err)
	} else {
		log.Printf("Put/update GameState associated with %s in DynamoDB table: GameStates", uuid)
	}

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

	all_rolls := []string{}
	for _, tr := range move_log {
		roll_strings := []string{"Player: " + tr.PlayerId}

		player_name, ok := id_to_name[tr.PlayerId]
		if !ok {
			log.Println("could not find player_name in id_to_name, player_name set to empty string")
		}

		for i := 0; i < 3; i++ {
			roll_strings = append(roll_strings, fmt.Sprintf("%s:%s : %s", player_name, tr.TurnResult[i][0], tr.TurnResult[i][1]))
		}

		if tr.IsDead == true || tr.ContinueTurn == false {

			turn_end_string := fmt.Sprintf("%s:%s, Total Score: %d, Turn Ended\n", player_name, tr.PlayerId, tr.TotalScore)
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

func send_chat_message(response http.ResponseWriter, request *http.Request) {
	log.Println("Sending Message")
	response.Header().Set("Content-type", "text/plain")

	err := request.ParseForm()
	if err != nil {
		log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	chat_id, err := parse_input(request, "chat_id")
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	player_id, err := parse_input(request, "player")
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	player_name, ok := id_to_name[player_id]
	if !ok {
		log.Println("could not find player_name in id_to_name, player_name set to empty string")
	}

	log.Print("PLAYERS NAME ID: ", player_id)

	log.Printf("chat id : %s", chat_id)
	zombie_chat, ok := zombie_chats[chat_id]

	if !ok {
		log.Printf("Could not find chat with id: %s\n", chat_id)
	}

	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
	if err != nil {
		fmt.Println("Could not parse request body")
	}

	err = request.Body.Close()
	if err != nil {
		fmt.Println("Could not close request Body")
	}

	message_info := map[string]string{}

	err = json.Unmarshal(body, &message_info)
	if err != nil {
		panic(err)
	}

	message := message_info["message"]

	player_message := fmt.Sprintf("%s:%s : %s", player_name, player_id, message)

	zombie_chat.ThreadSafeAppend(player_message)

	fmt.Fprintf(response, player_message)

}

func receive_all_chat_messages(response http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}

	chat_id_form, _ := request.Form["chat_id"]
	chat_id := chat_id_form[0]

	current_chat, ok := zombie_chats[chat_id]
	if !ok {
		log.Printf("Chat id has not been found")
		return
	}

	var messages_str string

	if len((*current_chat).Messages) >= 100 {
		messages_str = strings.Join((*current_chat).Messages[len((*current_chat).Messages)-10:], "\n")
	} else {
		messages_str = strings.Join((*current_chat).Messages, "\n")
	}

	fmt.Fprintf(response, messages_str)
}

func (i_to_n *id_to_name_type) set_player_name(response http.ResponseWriter, request *http.Request) {
	log.Println("Setting Name")

	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
	if err != nil {
		log.Println("Could not read Body")
	}

	err = request.Body.Close()
	if err != nil {
		log.Println("Could not close Body")
	}

	name_and_id_info := map[string]string{}

	err = json.Unmarshal(body, &name_and_id_info)
	if err != nil {
		log.Println("Could not Unmarshal Body into name_and_id_info")
	}

	player_name, ok := name_and_id_info["player_name"]
	if !ok {
		log.Println("player_name not in name_and_id_info")
	}

	player_id, ok := name_and_id_info["player_id"]
	if !ok {
		log.Println("player_id not in name_and_id_info")
	}

	log.Printf("Adding player_id, player_name: %s -> %s\n", player_id, player_name)
	(*i_to_n)[player_id] = player_name
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

	mux.HandleFunc("/zombie_dice_multi_player/send_chat_message", send_chat_message)
	mux.HandleFunc("/zombie_dice_multi_player/receive_all_chat_messages", receive_all_chat_messages)

	mux.HandleFunc("/zombie_dice_multi_player/set_player_name", id_to_name.set_player_name)

	mux.HandleFunc("/four_dice_roll", four_dice_roll)
	mux.HandleFunc("/roll_dice", roll_dice)

	log.Printf("Started dumb Dice web server! Try it on http://ip_address:8000")

	err := http.ListenAndServe("0.0.0.0:8000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
