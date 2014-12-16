package main

func n_dice_roll(n_sides int, n_rolls int) int {

	my_dice := InitDefaultDice(n_sides)
	my_sum := 0
	for i := 0; i < n_rolls; i++ {
		side := my_dice.Roll()
		my_sum += side.numerical_value
	}
	return my_sum
}
