package dice

import "testing"

func initBasicTestDeck(deck_size int) Deck {
	dice_list := make(Dices, deck_size)
	for i := 0; i < deck_size; i++ {
		dice_list[i] = InitDefaultDice(i + 1)
	}

	return Deck{Name: "Test Deck", Dices: dice_list}
}

func TestAllDiceCreated(t *testing.T) {
	deck := initBasicTestDeck(10)

	//Check that 10 dice were created
	if len(deck.Dices) != 10 {
		t.Errorf("Expected list of %d dice, but got %d",
			10, len(deck.Dices))
	}
}

func TestInitialDiceOrder(t *testing.T) {
	deck := initBasicTestDeck(10)

	//Ensure initial order is as expected
	for i := 0; i < 10; i++ {
		dice := deck.Dices[i]

		if len(dice.Sides) != i+1 {
			t.Errorf("Expected dice with %d sides, but got %d",
				i+1, len(dice.Sides))
		}
	}
}

func TestShuffle(t *testing.T) {
	//Making huge deck, since order of shuffled deck is tested
	//this minimizes the chance that we actually shuffle the deck,
	//but still get a perfectly sorted permutation
	const DICE_IN_DECK = 1000

	deck := initBasicTestDeck(DICE_IN_DECK)
	deck.Shuffle()

	//Ensure shuffled and all dice are present
	var dice_list_flag [DICE_IN_DECK]int  //Ensures all dice are still there
	var dice_list_order [DICE_IN_DECK]int //Tracks order
	for i := 0; i < DICE_IN_DECK; i++ {
		dice := deck.Dices[i]

		j := len(dice.Sides) - 1
		if dice_list_flag[j] != 0 {
			t.Errorf("Did not expect dice with same amount of sides after shuffling!")
		} else {
			dice_list_flag[j] = j
			dice_list_order[i] = j
		}
	}

	ordered := true
	for i := 0; i < DICE_IN_DECK; i++ {
		if i != dice_list_order[i] {
			ordered = false
			break
		}
	}

	if ordered {
		t.Errorf("Expected deck to be shuffled, but order remained the same!")
	}
}
