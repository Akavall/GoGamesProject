package zombie_dice

import (
	"github.com/Akavall/GoGamesProject/dice"
)

type ZombieDeck struct {
	Deck dice.Deck
}

func (d *ZombieDeck) Shuffle() {
	d.Deck.Shuffle()
}

func (d *ZombieDeck) DealDice(num_dice int) (dice.Dices, error) {
	return d.Deck.DealDice(num_dice)
}

func (d *ZombieDeck) AddDice(new_dice dice.Dice) {
	d.Deck.AddDice(new_dice)
}

func InitZombieDeck() ZombieDeck {

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

	zombie_dice_deck := dice.Deck{Name: "ZombieDiceDeck", Dices: dices}
	return ZombieDeck{Deck: zombie_dice_deck}
}

func make_slice_of_sides(string_sides []string) []dice.Side {
	sides := make([]dice.Side, len(string_sides))
	for ind, s := range string_sides {
		sides[ind] = dice.Side{Name: s}
	}
	return sides
}
