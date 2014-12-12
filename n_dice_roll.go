package main

import ("fmt")

func n_dice_roll(n_sides int, n_rolls int) int {
	my_dice := Dice{n_sides: n_sides}
        my_sum := 0
        i := 0 
	for i < n_rolls{
		my_sum += my_dice.Roll()
		i += 1
	}
	return my_sum
}

func main(){
	fmt.Println(n_dice_roll(6, 2))
}