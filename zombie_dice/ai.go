package zombie_dice

import (
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

func SimulationistAI(shots int, brains int, deck_left dice.Deck) int {
	n_iterations := 10000
	n_killed := 0
	n_brains := 0
	for i := 0; i < n_iterations; i++ {
		deck_left.Shuffle()
		n_inner_shots := 0
		n_inner_brains := 0
		for j := 0; j < 3; j++ {
			side := deck_left.Dices[j].Roll()
			if side.Name == "brain" {
				n_inner_brains++
			} else if side.Name == "shot" {
				n_inner_shots++
				if (n_inner_shots + shots) >= 3 {
					n_killed++
					continue
				}
				// If we did not get shot, we get the brains
				n_brains += n_inner_brains
			}
		}
	}

	// This commented out code illustrates the intentions
        // more clearly than simplified code below

	// expected_brains := float64(n_brains) / float64(n_iterations)
	// chance_to_get_killed := float64(n_killed) / float64(n_iterations)

	// if expected_brains > chance_to_get_killed * float64(brains) {
	// 	return 1
	// } else {
	// 	return 0
	// }

	if n_brains > n_killed * brains {
		return 1
	} else {
		return 0
	}

}
