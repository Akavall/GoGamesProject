package main

import (
	"fmt"
)

func n_dice_roll(n_sides int, n_rolls int) int {
	my_dice := Dice{n_sides: n_sides}
	my_sum := 0
	for i := 0; i < n_rolls; i++ {
		my_sum += my_dice.Roll()
	}
	return my_sum
}

func main() {
	fmt.Println(n_dice_roll(6, 2))
}
