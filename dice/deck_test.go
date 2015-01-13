package dice

import "testing"

const DEFAULT_DECK_SIZE = 10

func TestAllDiceCreated(t *testing.T) {
	deck := initBasicTestDeck(DEFAULT_DECK_SIZE)

	//Check that 10 dice were created
	if len(deck.Dices) != DEFAULT_DECK_SIZE {
		t.Errorf("Expected list of %d dice, but got %d",
			DEFAULT_DECK_SIZE, len(deck.Dices))
	}
}

func TestInitialDiceOrder(t *testing.T) {
	deck := initBasicTestDeck(DEFAULT_DECK_SIZE)

	//Ensure initial order is as expected
	if !isDeckOrdered(deck) {
		t.Errorf("Deck not created with expected ascending order of sides!")
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
	var dice_list_flag [DICE_IN_DECK]int
	for i := 0; i < DICE_IN_DECK; i++ {
		dice := deck.Dices[i]

		j := len(dice.Sides) - 1
		if dice_list_flag[j] != 0 {
			t.Errorf("Did not expect dice with same amount of sides after shuffling!")
		} else {
			dice_list_flag[j] = j
		}
	}

	if isDeckOrdered(deck) {
		t.Errorf("Expected deck to be shuffled, but order remained the same!")
	}
}

func TestDealDice(t *testing.T) {
	deck := initBasicTestDeck(DEFAULT_DECK_SIZE)
	dices, err := deck.DealDice(1)

	if err != nil {
		t.Errorf("Unexpected failure when dealing dice!")
	}

	if len(deck.Dices) != DEFAULT_DECK_SIZE-1 {
		t.Errorf(`Deck shrunk by more dice than dealt!
                Expected deck of size %d but got deck of size %d`,
			DEFAULT_DECK_SIZE-1, len(deck.Dices))
	}

	if len(dices) != 1 {
		t.Errorf(`Received more dice than expected! 
               Expected %d dice, but got %d`,
			1, len(dices))
	}
}

func TestDealDiceExact(t *testing.T) {
	// Deal exactly as many dice as there are in the deck
	deck := initBasicTestDeck(DEFAULT_DECK_SIZE)
	dices, err := deck.DealDice(DEFAULT_DECK_SIZE)

	if err != nil {
		t.Errorf("Unexpected failure when dealing dice!")
	}

	if len(dices) != DEFAULT_DECK_SIZE {
		t.Errorf("Expected entire deck to be dealt!")
	}

	if len(deck.Dices) != 0 {
		t.Errorf("Expected remaining deck to be zero!")
	}
}

func TestDealDiceFailure(t *testing.T) {
	deck := initBasicTestDeck(DEFAULT_DECK_SIZE)
	_, err := deck.DealDice(DEFAULT_DECK_SIZE + 1)

	if err == nil {
		t.Errorf(`Expected error after asking to deal more dice 
                than there are in the deck!`)
	}
}

func initBasicTestDeck(deck_size int) Deck {
	dice_list := make(Dices, deck_size)
	for i := 0; i < deck_size; i++ {
		dice_list[i] = InitDefaultDice(i + 1)
	}

	return Deck{Name: "Test Deck", Dices: dice_list}
}

func isDeckOrdered(deck Deck) bool {
	for i := 0; i < len(deck.Dices); i++ {
		dice := deck.Dices[i]

		if len(dice.Sides) != i+1 {
			return false
		}
	}

	return true
}
