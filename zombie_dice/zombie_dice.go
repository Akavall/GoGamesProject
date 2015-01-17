package zombie_dice

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	
	"github.com/Akavall/GoGamesProject/dice"
)

func initialize_deck() dice.Deck {

	green := []string{"shot", "walk", "walk", "brain", "brain", "brain"}
	yellow := []string{"shot", "shot", "walk", "walk", "brain", "brain"}
	red := []string{"shot", "shot", "shot", "walk", "walk", "brain"}

	green_sides := make_slice_of_sides(green)
	yellow_sides := make_slice_of_sides(yellow)
	red_sides := make_slice_of_sides(red)

	// Put dices in the deck

	const N_GREEN, N_YELLOW, N_RED = 6, 4, 3

	dices := make([]dice.Dice, 0)

	for i := 0; i < N_GREEN; i++ {
		dices = append(dices, dice.Dice{Name: "green", Sides: green_sides})
	}

	for i := 0; i < N_YELLOW; i++ {
		dices = append(dices, dice.Dice{Name: "yellow", Sides: yellow_sides})
	}

	for i := 0; i < N_RED; i++ {
		dices = append(dices, dice.Dice{Name: "red", Sides: red_sides})
	}

	zombie_dice_deck := dice.Deck{Name: "ZombieDiceDeck", Dices:dices}

	return zombie_dice_deck
}

func make_slice_of_sides(string_sides []string) []dice.Side {
	sides := make([]dice.Side, len(string_sides))
	for ind, s := range string_sides {
		sides[ind] = dice.Side{Name: s}
	}
	return sides
}

func players_turn(deck dice.Deck, player_name string) (int, error) {

	brains := 0
	shots := 0
	
	for {
		if len(deck.Dices) < 3 {
			fmt.Println("You have ran out of dices")
			fmt.Printf("Your final score is : %d", brains)
			return brains, nil
		}
		dices_to_roll, err := deck.DealDice(3)
		if err != nil {
			return 0, err
		}

		for _, d := range dices_to_roll {
			inner_walks := 0
			side := d.Roll()
			fmt.Println("You Rolled : ", d.Name, side.Name)
			if side.Name == "brain" {
				brains++
			} else if side.Name == "shot" {
				shots++
			} else {
				inner_walks++
			}
		}

		if shots >= 3 {
			fmt.Println("You have been shot 3 times, you've scored 0")
			return 0, nil
		}

		fmt.Printf("Your current score is %d\n", brains)
		fmt.Printf("Your have been shot %d times\n", shots)
		fmt.Println("Do you want to continue? Hit 1 to contintue and 0 to stop")

		var answer int

		switch player_name {
		case "human":
			answer = get_terminal_input()
		case "greedy":
			answer = GreedyAI(shots)
		case "careful":
                        answer = CarefulAI(shots)
                case "random":
			answer = RandomAI()
		
		}

		if answer == 0 {
			fmt.Println("You scored : ", brains)
			return brains, nil
		}
	}
	fmt.Println("The turn has ended")
	return brains, nil
}

func PlayWithAI() {
	player_total_score := 0
	ai_total_score := 0
	deck := initialize_deck()

	ai_name := select_ai()

	round_counter := 0
	for {
		round_counter++
		deck.Shuffle()
		player_score, err := players_turn(deck, "human")
		if err != nil {
			fmt.Println("Error Occurred on players turn")
			return
		}
		player_total_score += player_score

		fmt.Printf("Your total score is : %d\n", player_total_score)

		deck.Shuffle()
		ai_score, err_ai := players_turn(deck, ai_name)
		if err_ai != nil {
			fmt.Println("Error Occurred on ai turn")
			return
		}
			
		ai_total_score += ai_score

		fmt.Printf("Round : %d\n", round_counter)
		fmt.Printf("Your total score is : %d\n", player_total_score)
		fmt.Printf("AI total score is : %d\n", ai_total_score)

		if player_total_score >= 13 || ai_total_score >= 13 {
			if player_total_score > ai_total_score {
				fmt.Println("Congratulations You Won!")
				return
			} else if player_total_score < ai_total_score {
				fmt.Println("AI won! Better Luck Next Time!")
				return
			}
		}
	}
}

func get_terminal_input() int {
	reader := bufio.NewReader(os.Stdin)
	raw_string, _ := reader.ReadString('\n')
	clean_string := strings.Replace(raw_string, "\n", "", -1)
	answer, _ := strconv.Atoi(clean_string)
	return answer
}

func select_ai() string {
back:
	fmt.Println("Please Select an AI you want to play against")
	fmt.Printf("Greedy : press %d\n", 1)
	fmt.Printf("Careful : press %d\n", 2)
	fmt.Printf("Random : press %d\n", 3)

        answer := get_terminal_input()
	switch answer {
	case 1:
		return "greedy"
	case 2:
		return "careful"
	case 3:
		return "random"
	default:
		fmt.Println("This is not a valid selction, please try again")
		goto back
	}
}