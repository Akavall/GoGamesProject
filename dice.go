package main 

import ("fmt"
        "math/rand"
        "time")  

type Dice struct {
	n_sides int
}

func (d *Dice) roll() int {
	rand.Seed(time.Now().UTC().UnixNano())
        return rand.Int() % d.n_sides + 1
}

func main(){
	my_dice := Dice{n_sides: 6}
        fmt.Println(my_dice.roll())
}