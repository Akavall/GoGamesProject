package statistics

import (
	"math"
)

func CalcRollProbabilities(score int, n_dice int, n_sides int) (float64, float64, float64) {
	prob_of_roll := CalcRollProbability(score, n_dice, n_sides)

	prob_lower := 0.0
	for s := n_dice; s < score; s++ {
		prob_lower += CalcRollProbability(s, n_dice, n_sides)
	}

	prob_higher := 0.0
	for s := score + 1; s <= (n_dice * n_sides); s++ {
		prob_higher += CalcRollProbability(s, n_dice, n_sides)
	}
	return prob_of_roll, prob_lower, prob_higher
}

func CalcRollProbability(score int, n_dice int, n_sides int) float64 {
	n_perms := CalcRollPermutations(score, n_dice, n_sides)
	p_single := math.Pow(1.0/float64(n_sides), float64(n_dice))
	return float64(n_perms) * p_single
}

func CalcRollPermutations(score int, n_dice int, n_sides int) int {
	pc := PermCalculator{count: 0}
	return pc.calc_perms(n_dice, 0, n_sides, score)
}

type PermCalculator struct {
	count int
}

func (pc *PermCalculator) calc_perms(n_dice int, score int, n_sides int, target int) int {
	if n_dice == 1 {
		if target-score >= 1 && target-score <= n_sides {
			pc.count++
			return 1
		} else {
			return 0
		}
	}
	for i := 1; i <= n_sides; i++ {
		pc.calc_perms(n_dice-1, score+i, n_sides, target)
	}
	return pc.count
}
