package zombie_dice

import (
	"errors"
	"fmt"
	"log"

	"github.com/Akavall/GoGamesProject/dice"
)

//TO-DO: should probably define these to be configurable for each new game...
const WINNING_SCORE = 13
const DICE_TO_DEAL = 3
const SHOTS_UNTIL_DEAD = 3

type GameState struct {
	GameStateID string `dynamodbav:"game_state_id"`
	Players
	ZombieDeck
	PlayerTurn int
	Winner     Player
	GameOver   bool
	IsActive   bool
	MoveLog    []PlayerTurnResult
}

type Players []Player

type Player struct {
	*PlayerState
	Id         string
	IsAI       bool
	TotalScore *int
}

type PlayerState struct {
	TurnsTaken         int
	CurrentScore       int
	TimesShot          int
	BrainsRolled       int
	WalksTakenLastRoll int
	IsDead             bool
}

type PlayerTurnResult struct {
	TurnResult   [][]string
	RoundScore   int
	TimesShot    int
	TotalScore   int
	IsDead       bool
	Winner       string
	PlayerId     string
	ContinueTurn bool
}

func (gs *GameState) EndTurn() {
	next_player_turn := gs.PlayerTurn + 1

	if next_player_turn >= len(gs.Players) {
		next_player_turn = 0
		gs.endRound()
	}

	gs.PlayerTurn = next_player_turn

	deck := InitZombieDeck()
	deck.Shuffle()
	gs.ZombieDeck = deck
}

func (gs *GameState) endRound() {

	player_score_to_count := map[int]int{}
	max_score := 0
	for _, p := range gs.Players {
		log.Printf("\033[0;32mPlayer: %s, Score: %d\033[0m", p.Id, *p.TotalScore)
		player_score_to_count[*p.TotalScore] += 1
		if *p.TotalScore > max_score {
			max_score = *p.TotalScore
		}
	}

	log.Printf("\033[32mMax score: %d\033[0m", max_score)

	if max_score >= WINNING_SCORE && player_score_to_count[max_score] == 1 {
		for _, p := range gs.Players {
			if *p.TotalScore == max_score {
				gs.Winner = p
				gs.GameOver = true
			}
		}
	}

}

func (ps *PlayerState) Reset() {
	ps.TurnsTaken = 0
	ps.CurrentScore = 0
	ps.TimesShot = 0
	ps.WalksTakenLastRoll = 0
	ps.BrainsRolled = 0
	ps.IsDead = false
}

func InitGameState(players Players, gameStateId string) (gs GameState, err error) {
	deck := InitZombieDeck()
	deck.Shuffle()

	return GameState{GameStateID: gameStateId, Players: players, ZombieDeck: deck, PlayerTurn: 0, Winner: Player{}, GameOver: false, IsActive: false}, nil
}

func (p *Player) TakeTurn(deck *ZombieDeck) (s [][]string, err error) {
	// turn_result := [3][2]string{}

	turn_result := make([][]string, 3)

	for i := 0; i < 3; i++ {
		turn_result[i] = make([]string, 2)
	}

	log.Printf("p address: %p", &p)
	log.Printf("In TakeTurn is dead: %t", p.PlayerState.IsDead)

	if p.PlayerState.IsDead == true {
		return turn_result, errors.New(fmt.Sprintf("Player %s is dead and cannot take more turns!", p.Id))
	}

	if len(deck.Deck.Dices) < DICE_TO_DEAL {
		log.Printf("\033[33mDeck size is too small: %d, adding another zombie deck to the existing deck\033[0m", len(deck.Deck.Dices))
		new_deck := InitZombieDeck()
		new_deck.Shuffle()
		deck.Deck.PrependDeck(new_deck.Deck)
	}

	dices_to_roll, err := deck.DealDice(DICE_TO_DEAL)
	if err != nil {
		return
	}

	p.PlayerState.WalksTakenLastRoll = 0
	sides := make([]dice.Side, 0)
	for roll_ind, d := range dices_to_roll {
		side := d.Roll()
		sides = append(sides, side)
		turn_result[roll_ind][0] = d.Name
		turn_result[roll_ind][1] = side.Name
		log.Printf("%s rolled: %s, %s\n", p.Id, d.Name, side.Name)

		if side.Name == "brain" {
			p.PlayerState.CurrentScore++
			p.PlayerState.BrainsRolled++
		} else if side.Name == "shot" {
			p.PlayerState.TimesShot++
		} else if side.Name == "walk" {
			// Since walks get replayed we have to
			// put them back in the deck
			deck.AddDice(d)
			p.PlayerState.WalksTakenLastRoll++
		} else {
			return turn_result, errors.New(fmt.Sprintf("Unrecognized dice side has been rolled: %s", side.Name))
		}
	}

	if p.PlayerState.TimesShot >= SHOTS_UNTIL_DEAD {
		p.PlayerState.IsDead = true
		log.Printf("p address: %p", &p)
		log.Printf("%s is dead", p.Id)
	}

	p.PlayerState.TurnsTaken++
	return turn_result, nil //TO-DO: need proper return here that significies dice color + side rolled
}

func InitPlayerState() *PlayerState {
	return &PlayerState{TurnsTaken: 0, CurrentScore: 0, TimesShot: 0, IsDead: false}
}
