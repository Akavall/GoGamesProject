package zombie_dice

import (
	"fmt"
	"github.com/Akavall/GoGamesProject/dice"
)

func GreedyAI(shots int) int {
	if shots == 2 {
		return 0
	}
	return 1
}

func CarefulAI(shots int) int {
	if shots < 2 {
		return 0
	}
	return 1
}

func RandomAI() int {
	two_sided_dice := dice.InitDefaultDice(2)
	return two_sided_dice.Roll().Numerical_value - 1
}

<<<<<<< HEAD
func SimulationistAI(shots, brains, walks int, deck_left dice.Deck) int {
	// Refactor this:
=======
func SimulationistAI(previous_shots, already_gained_brains, walks int, deck_left dice.Deck) int {
	// This is a dumb simulationist it misses walk dices
>>>>>>> 72d4c1c3c4f1984d71c0ae0f628471aa37b122a9
	n_iterations := 10000
	all_killed := 0
	all_brains := 0
	walk_dices, err := deck_left.DealDice(walks)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < n_iterations; i++ {
		deck_left.Shuffle()
		
		n_killed, n_brains := simulate_one_roll(&walk_dices, &deck_left, previous_shots)
		all_killed += n_killed
		all_brains += n_brains
		
	}

	fmt.Println("all_brains", all_brains)
	fmt.Println("to lose", all_killed * already_gained_brains)

	if all_brains > all_killed * already_gained_brains {
		return 1
	} else {
		return 0
	}

}

func simulate_one_roll(walk_dices *dice.Dices, deck_left *dice.Deck, previous_shots int) (int, int) {
	n_inner_shots := 0
	n_inner_brains := 0
	for j := 0; j < 3; j++ {
		var side dice.Side 
		if j < len(*walk_dices) {
			side = (*walk_dices)[j].Roll()
		} else {                        
			side = deck_left.Dices[j].Roll()
		}

		if side.Name == "brain" {
			n_inner_brains++
		} else if side.Name == "shot" {
			n_inner_shots++
			if (n_inner_shots + previous_shots) >= 3 {
				return 1, 0
			}
		}
	}
	return 0, n_inner_brains

}
