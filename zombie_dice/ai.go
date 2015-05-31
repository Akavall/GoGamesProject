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

func SimulationistAI(previous_shots, already_gained_brains, walks int, deck_left *ZombieDeck) int {

	// We need to make a copy of deck_left, so not to
	// mutate the original deck when we are training 

	dices_copy := make([]dice.Dice, len(deck_left.Deck.Dices))
	copy(dices_copy, deck_left.Deck.Dices)

	deck_copy := dice.Deck{Name: "deck_copy", Dices: dices_copy}
	zombie_deck_c := ZombieDeck{Deck: deck_copy}

	n_iterations := 10000
	all_killed := 0
	all_brains := 0
	walk_dices, err := zombie_deck_c.DealDice(walks)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < n_iterations; i++ {
		zombie_deck_c.Shuffle()

		n_killed, n_brains := simulate_one_roll(&walk_dices, &zombie_deck_c, previous_shots)
		all_killed += n_killed
		all_brains += n_brains

	}

	fmt.Println("all_brains", all_brains)
	fmt.Println("to lose", all_killed*already_gained_brains)

	if all_brains > all_killed*already_gained_brains {
		return 1
	} else {
		return 0
	}

}

func simulate_one_roll(walk_dices *dice.Dices, zombie_deck_c *ZombieDeck, previous_shots int) (int, int) {
	n_inner_shots := 0
	n_inner_brains := 0
	for j := 0; j < 3; j++ {
		var side dice.Side
		if j < len(*walk_dices) {
			side = (*walk_dices)[j].Roll()
		} else {
			side = zombie_deck_c.Deck.Dices[j-len(*walk_dices)].Roll()
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
