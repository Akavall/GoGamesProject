package main

import (
	"math/rand"
	"time"
)

type Dice struct {
	n_sides int
}

func (d *Dice) Roll() int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int()%d.n_sides + 1
}
