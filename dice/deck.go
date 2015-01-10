package dice

import (
	"math/rand"
	"time"
)

type Deck struct {
	Name string
	Dices
}

type Dices []Dice

func (d *Deck) Shuffle() {
	if len(d.Dices) <= 1 {
		return
	}

	rand.Seed(time.Now().UTC().UnixNano())
	rand_inds := rand.Perm(len(d.Dices))

	shuffled_dice := make(Dices, len(d.Dices))
	for i, rand_ind := range rand_inds {
		shuffled_dice[i] = d.Dices[rand_ind]
	}

	d.Dices = shuffled_dice
}
