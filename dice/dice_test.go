package dice

import "testing"

func TestInitDefaultDice(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{1, 1},
		{6, 6},
		{10, 10},
	}

	for _, c := range cases {
		dice := InitDefaultDice(c.in)
		if len(dice.Sides) != c.want {
			t.Errorf("Expected %d side(s), but got %d", c.in, len(dice.Sides))
		}

		for i := 0; i < len(dice.Sides); i++ {
			if dice.Sides[i].Numerical_value != i+1 {
				t.Errorf("Expected dice side with value %d, but got %d",
					i+1, dice.Sides[i].Numerical_value)
			}
		}
	}
}
