package dice

import (
	"math/rand"
	"time"
)

type Dice struct {
	Name        string
	Description string
	Sides
}

type Dices []Dice

func (d Dice) Roll() Side {
	rand.Seed(time.Now().UTC().UnixNano())
	random_roll := rand.Int() % len(d.Sides)
	return d.Sides[random_roll]
}

func (d Dice) RollNTimes(n_rolls int) Sides {
	Sides := make([]Side, n_rolls)
	for i := 0; i < n_rolls; i++ {
		Sides[i] = d.Roll()
	}
	return Sides
}

func InitDefaultDice(n_sides int) Dice {
	Sides := make([]Side, n_sides)
	for i := 0; i < n_sides; i++ {
		Sides[i] = Side{Numerical_value: i + 1}
	}
	return Dice{Name: "Default dice", Sides: Sides}
}
