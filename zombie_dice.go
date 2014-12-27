package main

import (
	"fmt"
	"time"
	"math/rand"
	"bufio"
	"os"
	"strconv"
	"strings"

        "github.com/Akavall/GoGamesProject/dice"
)

func initialize_deck () []dice.Dice {

	green := []string{"shot", "walk", "walk", "brain", "brain", "brain"}
	yellow := []string{"shot", "shot", "walk", "walk", "brain", "brain"}
	red := []string{"shot", "shot", "shot", "walk", "walk", "brain"}
        
	green_sides := make_slice_of_sides(green)
	yellow_sides := make_slice_of_sides(yellow)
        red_sides := make_slice_of_sides(red)
	
        // Put dices in the deck

        const N_GREEN, N_YELLOW, N_RED = 6, 4, 3

	deck := make([]dice.Dice, 0)

	for i := 0; i < N_GREEN; i++ {
		deck = append(deck, dice.Dice{Name: "green", Sides: green_sides})
	}

	for i := 0; i < N_YELLOW; i++ {
		deck = append(deck, dice.Dice{Name: "yellow", Sides: yellow_sides})
	}

	for i := 0; i < N_RED; i++ {
		deck = append(deck, dice.Dice{Name: "red", Sides: red_sides})
	}

	// Shuffle the deck
	// rand.Seed(time.Now().UTC().UnixNano())

	// rand_inds := rand.Perm(N_DICES)
        // shuffled_deck := make([]int, N_DICES)
  
        // for ind, rand_ind := range rand_inds {
	// 	shuffled_deck[ind] = deck[rand_ind]
	// }
	
	return deck
}

func shuffle_deck(deck []dice.Dice) []dice.Dice {
	rand.Seed(time.Now().UTC().UnixNano())
	rand_inds := rand.Perm(len(deck))
	shuffled_deck := make([]dice.Dice, len(deck))

	for ind, rand_ind := range rand_inds {
		shuffled_deck[ind] = deck[rand_ind]
	}
	return shuffled_deck
}

func make_slice_of_sides(string_sides []string) []dice.Side {
	sides := make([]dice.Side, len(string_sides))
	for ind, s := range string_sides {
		sides[ind] = dice.Side{Name: s}
	}
	return sides
}

func players_turn(deck []dice.Dice) int {
	// Need to do major refactoring here
	// deck = initialize_deck()
	brains := 0
        shots := 0
	old_dices := make([]dice.Dice, 0)
        reader := bufio.NewReader(os.Stdin)

	// While loop
	for i := 0; i < 1; i += 0 {
		if (len(deck) + len(old_dices) < 3) {
			fmt.Println("You have ran out of dices")
                        fmt.Printf("Your final score is : %d", brains)
			return brains
		}
		dices_to_roll := pop_last_n(&deck, 3 - len(old_dices))
		dices_to_roll = append(dices_to_roll, old_dices...)
		old_dices = nil
		for _, d := range dices_to_roll {
			inner_walks := 0
			side := d.Roll()
			fmt.Println("You Rolled : ", d.Name, side.Name)
			if (side.Name == "brain") {
				brains++
			} else if (side.Name == "shot") {
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

		fmt.Printf("Your current score is %d\n", brains)
		fmt.Printf("Your have been shot %d times\n", shots)
		fmt.Println("Do you want to continue? Hit 1 to contintue and 0 to stop")
		raw_string, _ := reader.ReadString('\n')
		clean_string := strings.Replace(raw_string, "\n", "", -1)
		answer, _ := strconv.Atoi(clean_string)

		if (answer == 0) {
			fmt.Println("You scored : ", brains)
			return brains
		}
		
		// emptying the slice
	}
	fmt.Println("The turn has ended")
	return brains
}

func pop_last_n(a_ptr *[]dice.Dice, n_to_pop int) []dice.Dice {

	a := *a_ptr
	poped_slice := a[len(a) - n_to_pop : len(a)]
	a = append(a[:0], a[:len(a) - n_to_pop]...)
        *a_ptr = a
        
	return poped_slice
}

// func monster_dice() {
// }

func main() {
	deck := initialize_deck()
	shuffled_deck := shuffle_deck(deck)
	players_turn(shuffled_deck)
}
