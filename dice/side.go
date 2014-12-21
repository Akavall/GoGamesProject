package dice

type Side struct {
	Name            string
	Description     string
	Numerical_value int
}

type Sides []Side

func (s Sides) SumSides() int {
	sum := 0
	for i := 0; i < len(s); i++ {
		sum += s[i].Numerical_value
	}
	return sum
}
