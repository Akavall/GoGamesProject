package zombie_dice

import (
	"testing"
	"github.com/Akavall/GoGamesProject/dice"
)


func make_green_dice() dice.Dice {

	green_sides := make([]dice.Side, 6)
	green_sides[0] = dice.Side{Name: "brain"}
	green_sides[1] = dice.Side{Name: "brain"}
	green_sides[2] = dice.Side{Name: "brain"}
	green_sides[3] = dice.Side{Name: "shot"}
	green_sides[4] = dice.Side{Name: "walk"}
	green_sides[5] = dice.Side{Name: "walk"}
	green_dice := dice.Dice{Name: "green", Sides: green_sides}

	return green_dice
}

func make_yellow_dice() dice.Dice {

	yellow_sides := make([]dice.Side, 6)
	yellow_sides[0] = dice.Side{Name: "brain"}
	yellow_sides[1] = dice.Side{Name: "brain"}
	yellow_sides[2] = dice.Side{Name: "shot"}
	yellow_sides[3] = dice.Side{Name: "shot"}
	yellow_sides[4] = dice.Side{Name: "walk"}
	yellow_sides[5] = dice.Side{Name: "walk"}
	yellow_dice := dice.Dice{Name: "yellow", Sides: yellow_sides}

	return yellow_dice
}

func make_red_dice() dice.Dice {

	red_sides := make([]dice.Side, 6)
	red_sides[0] = dice.Side{Name: "brain"}
	red_sides[1] = dice.Side{Name: "shot"}
	red_sides[2] = dice.Side{Name: "shot"}
	red_sides[3] = dice.Side{Name: "shot"}
	red_sides[4] = dice.Side{Name: "walk"}
	red_sides[5] = dice.Side{Name: "walk"}
	red_dice := dice.Dice{Name: "red", Sides: red_sides}

	return red_dice
}

func TestSimulateOneRoll(t *testing.T) {

	walk_dices_1 := make(dice.Dices, 3)
	walk_dices_1[0] = make_green_dice()
	walk_dices_1[1] = make_green_dice()
	walk_dices_1[2] = make_green_dice()

	walk_dices_2 := make(dice.Dices, 3)
	walk_dices_2[0] = make_red_dice()
	walk_dices_2[1] = make_red_dice()
	walk_dices_2[2] = make_red_dice()

	deck_left := dice.Deck{Dices: make(dice.Dices, 0)}

	killed_1 := 0
	killed_2 := 0
	brains_1 := 0
	brains_2 := 0

	for i := 0; i < 1000; i++ {

		k_1, b_1 := simulate_one_roll(&walk_dices_1, &deck_left, 0)
		k_2, b_2 := simulate_one_roll(&walk_dices_2, &deck_left, 0)

		killed_1 += k_1
		killed_2 += k_2
		brains_1 += b_1
		brains_2 += b_2

	}

	if killed_1 >= killed_2 {
		t.Errorf("Expected green dice walks to yield less kills than red dice walks")
	}

	if brains_1 <= brains_2 {
		t.Errorf("Expected greed dice walks to yield more kills than red dice walks")
	}

}

