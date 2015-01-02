package zombie_dice

import ( "testing"

	"github.com/Akavall/GoGamesProject/dice"
)

func TestInitializeDeck(t *testing.T) {
	deck := initialize_deck()

	expected_names := map[int]string{0: "green",
		1:  "green",
		2:  "green",
		3:  "green",
		4:  "green",
		5:  "green",
		6:  "yellow",
		7:  "yellow",
		8:  "yellow",
		9:  "yellow",
		10: "red",
		11: "red",
		12: "red"}

	for ind, d := range deck {
		if d.Name != expected_names[ind] {
			t.Errorf("Expected %d: %s, but got %s", ind, expected_names[ind], d.Name)
		}
	}

}

func TestShuffleDeck(t *testing.T) {
	green := dice.Dice{Name: "green", Sides: make([]dice.Side, 0)}
	yellow := dice.Dice{Name: "yellow", Sides: make([]dice.Side, 0)}
	red := dice.Dice{Name: "red", Sides: make([]dice.Side, 0)}

	deck := []dice.Dice{green, yellow, red}

	first_position := make([]string, 100)

	for i := 0; i < 100; i++ {
		shuffled_deck := shuffle_deck(deck)
		inner_green := 0
		inner_yellow := 0
		inner_red := 0
		for ind, d := range shuffled_deck {
			if d.Name == "green" {
				inner_green++
			} else if d.Name == "yellow" {
				inner_yellow++
			} else if d.Name == "red" {
				inner_red++
			}
			if ind == 0 {
				first_position[i] = d.Name
			}
		}
		if inner_green*inner_red*inner_yellow != 1 {
			t.Errorf("Got %d green, %d red, %d yellow, while expecting one of each", inner_green, inner_red, inner_yellow)
		}

	}

	counter := map[string]int{"green": 0, "red": 0, "yellow": 0}
	for _, name := range first_position {
		counter[name]++
	}

	for key, value := range counter {
		if value == 0 {
			t.Errorf("name : %s showed up in 1st position 0 times", key)
		}
	}
}

