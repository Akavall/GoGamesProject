package dice 

import (
	"math/rand"
	"time"
)

type Side struct {
	name            string
	description     string
	Numerical_value int
}

type Sides []Side

type Dice struct {
	name        string
	description string
	Sides
}

func (d *Dice) Roll() Side {
	rand.Seed(time.Now().UTC().UnixNano())
	random_roll := rand.Int() % len(d.Sides)
	return d.Sides[random_roll]
}

func InitDefaultDice(n_sides int) Dice {
	Sides := make([]Side, n_sides)
	for i := 0; i < n_sides; i++ {
		Sides[i] = Side{Numerical_value: i + 1}
	}
	return Dice{name: "Default dice", Sides: Sides}
}
