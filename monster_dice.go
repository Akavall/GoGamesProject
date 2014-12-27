package main

import (
	"time"
	"math/rand"
	"bufio"
	"os"
	"strconv"
	"strings"x

	"github.com/Akavall/GoGamesProject/dice"
)

func initialize_deck () []Dice {

	green_sides := []string{"shot", "walk", "walk", "brain", "brain", "brain"}
	yellow_sides := []string{"shot", "shot", "walk", "walk", "brain", "brain"}
	red_sides := []string{"shot", "shot", "shot", "walk", "walk", "brain"}
        
        // Put dices in the deck

        const N_DICES := 13
        const N_GREEN, N_YELLOW, N_RED := 6, 4, 3

	deck = make([]Dice, N_DICES)

	for i := 0; i < N_GREEN; i++ {
		deck = append(deck.dices, Dice{Name: "green", Sides: green_sides})
	}

	for i := 0; i < N_YELLOW; i++ {
		deck = append(deck.dices, Dice{Name: "yellow", Sides: yellow_sides})
	}

	for i := 0; i < N_RED; i++ {
		deck = append(deck.dices, Dice{Name: "red", Sides: red_sides})
	}

	// Shuffle the deck
	rand.Seed(time.Now().UTC().UnixNano())

	rand_inds = rand.Perm(N_DICES)
        shuffled_deck = make([]int, N_DICES)
  
        for ind, rand_ind := range rand_inds {
		suffled_deck[ind] = deck[rand_ind]
	}
	
	return shuffled_deck
}

func players_go(deck []Dice) int {
	// Need to do major refactoring here
	deck = initialize_deck()
	brains := 0
        shots := 0
	old_dices = int[]{}

	// While loop
	for i := 0; i < 1; i++ {
		if (len(deck) + len(old_dices) < 3) {
			fmt.Println("You have ran out of dices")
			return brains
		}
		dices_to_roll := pop_last_n(&deck, 3 - len(old_dices))
		dices_to_roll = append(dices_to_roll, old_dices)
		for _, d := range dices_to_roll {
			inner_walks := 0
			roll := d.Roll()
			fmt.Println("You Rolled : ", d.Name, roll)
			if (roll == "brain") {
				brains++
			} else if (roll == "shot") {
				shots++
			} else {
				inner_walks++
			old_dices = append(old_dices, d)
			}
		}

		if (shots >= 3) {
			fmt.Println("You have been shot 3 times, you've scored 0")
			return 0
		}

		fmt.Println("Do you want to continue? Hit 1 to contintue and 0 to stop")
		raw_string, _ := reader.ReadString('\n')
		clean_string := strings.Replace(raw_string, "\n", "", -1)
		answer, _ := strconv.Atoi(clean_string)

		if (answer == 0) {
			fmt.Println("You scored : ", brains)
			return brains
		}
		
		old_dices = int[]{}
	}
}

func pop_last_n(a_ptr *[]Dice, n_to_pop int) []Dice {

	a := *a_ptr
	poped_slice := a[len(a) - n_to_pop : len(a)]
	a = append(a[:0], a[:len(a) - n_to_pop]...)
        *a_ptr = a
        
	return poped_slice
}

// func monster_dice() {
// }
