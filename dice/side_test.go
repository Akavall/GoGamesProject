package dice

import "testing"

func TestSumSides(t *testing.T) {
	//When dice with 1 side with value 1 is rolled 10 times
	//Expect final sum to be 10
	sum := InitDefaultDice(1).RollNTimes(10).SumSides()

	if sum != 10 {
		t.Errorf("Expected sum %d, but got %d", 10, sum)
	}

	//When dice with 3 sides all with value 3 is rolled 5 times
	//Expect final sum to be 15
	Sides := make([]Side, 3)
	Sides[0] = Side{Numerical_value: 3}
	Sides[1] = Side{Numerical_value: 3}
	Sides[2] = Side{Numerical_value: 3}
	dice := Dice{Sides: Sides}

	sum = dice.RollNTimes(5).SumSides()
	if sum != 15 {
		t.Errorf("Expected sum %d, but got %d", 10, sum)
	}
}
