package dice

import (
	"errors"
	"math/rand"
	"time"
)

type Deck struct {
	Name string
	Dices
}

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

func (d *Deck) DealDice(num_dice int) (Dices, error) {
	if num_dice <= len(d.Dices) {
		dealt_dice := d.Dices[len(d.Dices)-num_dice : len(d.Dices)]
		d.Dices = append(d.Dices[:0], d.Dices[:len(d.Dices)-num_dice]...)
		return dealt_dice, nil
	} else {
		return nil, errors.New("Not enough dice in deck!")
	}
}

func (d *Deck) AppendDice(new_dice Dice) {
	d.Dices = append(d.Dices, new_dice)
}
