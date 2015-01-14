package dice_war

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Akavall/GoGamesProject/dice"
)

func DiceWar(n_your_dice_sides, n_your_dice, n_my_dice_sides, n_my_dice int) (int, int, string) {
	your_dice := dice.InitDefaultDice(n_your_dice_sides)
	my_dice := dice.InitDefaultDice(n_my_dice_sides)

	your_roll_sum := your_dice.RollNTimes(n_your_dice).SumSides()
	my_roll_sum := my_dice.RollNTimes(n_my_dice).SumSides()

	var result_string string

	if your_roll_sum > my_roll_sum {
		result_string = "You win! Well done"
	} else if your_roll_sum < my_roll_sum {
		result_string = "I win!, HaHa!"
	} else {
		result_string = "We are evenly matched, good game!"
	}
	return your_roll_sum, my_roll_sum, result_string
}

func DiceWarForTerminal() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ready for dice war?")

	fmt.Println("How many dice should you have?")
	raw_string, _ := reader.ReadString('\n')
	clean_string := strings.Replace(raw_string, "\n", "", -1)
	n_your_dice, _ := strconv.Atoi(clean_string)

	fmt.Println("How sides should your dice have?")
	raw_string, _ = reader.ReadString('\n')
	clean_string = strings.Replace(raw_string, "\n", "", -1)
	n_your_dice_sides, _ := strconv.Atoi(clean_string)

	fmt.Println("How many dice should I have?")
	raw_string, _ = reader.ReadString('\n')
	clean_string = strings.Replace(raw_string, "\n", "", -1)
	n_my_dice, _ := strconv.Atoi(clean_string)

	fmt.Println("How sides should my dice have?")
	raw_string, _ = reader.ReadString('\n')
	clean_string = strings.Replace(raw_string, "\n", "", -1)
	n_my_dice_sides, _ := strconv.Atoi(clean_string)

	your_roll_sum, my_roll_sum, result_string := DiceWar(n_your_dice_sides, n_your_dice, n_my_dice_sides, n_my_dice)

	fmt.Printf("You rolled : %d\n", your_roll_sum)
	fmt.Printf("I rolled : %d\n", my_roll_sum)
	fmt.Println(result_string)
}
