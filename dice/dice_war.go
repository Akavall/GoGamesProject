package dice 

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func dice_war() {
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

	your_roll_sum := N_dice_roll(n_your_dice_sides, n_your_dice)
	my_roll_sum := N_dice_roll(n_my_dice_sides, n_my_dice)

	fmt.Println("Your sum is : ", your_roll_sum)
	fmt.Println("My sum is : ", my_roll_sum)

	if your_roll_sum > my_roll_sum {
		fmt.Println("You win! Well done")
	} else if your_roll_sum < my_roll_sum {
		fmt.Println("I win!, HaHa!")
	} else {
		fmt.Println("We are evenly matched, good game!")
	}

}
