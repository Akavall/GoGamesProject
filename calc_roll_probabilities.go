package main

func CalcRollPermutations(score int, n_dice int, n_sides int) int {
	pc := PermCalculator{count: 0}
	return pc.calc_combs(n_dice, 0, n_sides, 0, score)
}

type PermCalculator struct {
	count int
}

func (pc *PermCalculator) calc_combs(n_dice int, score int, n_sides int, roll int, target int) int {
	if n_dice == 1 {
		if target-score >= 1 && target-score <= n_sides {
			pc.count++
			return 1
		} else {
			return 0
		}
	}
	for i := 1; i <= n_sides; i++ {
		pc.calc_combs(n_dice-1, score+i, n_sides, i, target)
	}
	return pc.count
}
