package zombie_dice

import (
	"testing"
)

func TestInitializeDeck(t *testing.T) {
	deck := InitZombieDeck()

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

	for ind, d := range deck.Deck.Dices {
		if d.Name != expected_names[ind] {
			t.Errorf("Expected %d: %s, but got %s", ind, expected_names[ind], d.Name)
		}
	}

}
