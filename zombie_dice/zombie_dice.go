package zombie_dice

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func players_turn(deck ZombieDeck, player_name string) (int, error) {

	brains := 0
	shots := 0

	for {
		if len(deck.Deck.Dices) < 3 {
			fmt.Printf("%s have ran out of dices", player_name)
			fmt.Printf("%s final score is : %d", player_name, brains)
			return brains, nil
		}

		dices_to_roll, err := deck.DealDice(3)
		if err != nil {
			return 0, err
		}

		for _, d := range dices_to_roll {
			side := d.Roll()
			fmt.Printf("%s rolled : %s, %s\n", player_name, d.Name, side.Name)
			if side.Name == "brain" {
				brains++
			} else if side.Name == "shot" {
				shots++
			} else {
				// Since walks get replayed we have to
				// put them back in the deck
				deck.AddDice(d)
			}
		}

		if shots >= 3 {
			fmt.Printf("%s have been shot 3 times, %s scored 0\n\n", player_name, player_name)
			fmt.Println("turn ending")
			time.Sleep(3 * 1e9)
			return 0, nil
		}

		fmt.Printf("%s current score is : %d\n", player_name, brains)
		fmt.Printf("%s have been shot : %d times\n", player_name, shots)
		fmt.Println("Do you want to continue? Hit 1 to contintue and 0 to stop")

		if player_name != "human" {
			time.Sleep(5 * 1e9)
		}

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
			fmt.Printf("%s %d : \n", player_name, brains)
			if player_name != "human" {
				fmt.Println("turn ending...")
				time.Sleep(3 * 1e9)
			}
			return brains, nil
		}
	}
	fmt.Println("The turn has ended")
	return brains, nil
}

func PlayWithAI() {
	player_total_score := 0
	ai_total_score := 0
	deck := InitZombieDeck()

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
		fmt.Printf("AI : %s total score is : %d\n", ai_name, ai_total_score)

		if player_total_score >= 13 || ai_total_score >= 13 {
			if player_total_score > ai_total_score {
				fmt.Println("Congratulations You Won!")
				return
			} else if player_total_score < ai_total_score {
				fmt.Printf("AI : %s won! Better Luck Next Time!\n", ai_name)
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
