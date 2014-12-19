package main

import (
	"testing"
)

func TestCalcRollPermutations_TwoDice_SixSides(t *testing.T) {

	const N_DICE, N_SIDES = 2, 6

	expected_score_to_n_pos := make(map[int]int)

	for i := 1; i <= N_SIDES; i++ {
		for j := 1; j <= N_SIDES; j++ {
			expected_score_to_n_pos[i+j] += 1
		}
	}

	for i := N_SIDES; i <= N_DICE*N_SIDES; i++ {
		result := CalcRollPermutations(i, N_DICE, N_SIDES)
		expected := expected_score_to_n_pos[i]
		if result != expected {
			t.Errorf("For score %d with %d dice and %d sides expected %d but got %d", i, N_DICE, N_SIDES, expected, result)
		}
	}
}

func TestCalcRollPermutations_ThreeDice_SevenSides(t *testing.T) {

	const N_DICE, N_SIDES = 3, 7

	expected_score_to_n_pos := make(map[int]int)

	for i := 1; i <= N_SIDES; i++ {
		for j := 1; j <= N_SIDES; j++ {
			for k := 1; k <= N_SIDES; k++ {
				expected_score_to_n_pos[i+j+k] += 1
			}
		}
	}

	for i := N_SIDES; i <= N_DICE*N_SIDES; i++ {
		result := CalcRollPermutations(i, N_DICE, N_SIDES)
		expected := expected_score_to_n_pos[i]
		if result != expected {
			t.Errorf("For score %d with %d dice and %d sides expected %d but got %d", i, N_DICE, N_SIDES, expected, result)
		}
	}
}

func TestCalcRollPermutations_FiveDice_TenSides(t *testing.T) {

	const N_DICE, N_SIDES = 5, 10

	expected_score_to_n_pos := make(map[int]int)

	for i := 1; i <= N_SIDES; i++ {
		for j := 1; j <= N_SIDES; j++ {
			for k := 1; k <= N_SIDES; k++ {
				for l := 1; l <= N_SIDES; l++ {
					for m := 1; m <= N_SIDES; m++ {
						expected_score_to_n_pos[i+j+k+l+m] += 1
					}
				}
			}
		}
	}

	for i := N_SIDES; i <= N_DICE*N_SIDES; i++ {
		result := CalcRollPermutations(i, N_DICE, N_SIDES)
		expected := expected_score_to_n_pos[i]
		if result != expected {
			t.Errorf("For score %d with %d dice and %d sides expected %d but got %d", i, N_DICE, N_SIDES, expected, result)
		}
	}
}
