package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCronIntervals_Empty(t *testing.T) {
	assert.Empty(t, GenerateCronIntervals("", "", "", "", ""))
}

func TestGenerateCronIntervals_Ranges(t *testing.T) {
	assert.Equal(
		t,
		[]map[string]int{
			{"Minute": 1},
			{"Minute": 2},
			{"Minute": 3},
			{"Minute": 5},
			{"Minute": 35},
		},
		GenerateCronIntervals("35,5,1-3", "", "", "", ""),
	)
}

func TestGenerateCronIntervals_StarsIgnored(t *testing.T) {
	assert.Empty(t, GenerateCronIntervals("*", "*", "*", "*", "*"))
}

func TestGenerateCronIntervals_StarDivisors(t *testing.T) {
	assert.Equal(
		t,
		[]map[string]int{
			{"Minute": 0},
			{"Minute": 5},
			{"Minute": 10},
			{"Minute": 15},
			{"Minute": 20},
			{"Minute": 25},
			{"Minute": 30},
			{"Minute": 35},
			{"Minute": 40},
			{"Minute": 45},
			{"Minute": 50},
			{"Minute": 55},
		},
		GenerateCronIntervals("*/5", "", "", "", ""),
	)
}

func TestGenerateCronIntervals_SingleCombo(t *testing.T) {
	assert.Equal(
		t,
		[]map[string]int{
			{"Minute": 0, "Hour": 1, "Day": 2, "Month": 3, "Weekday": 4},
		},
		GenerateCronIntervals("0", "1", "2", "3", "4"),
	)
}

func TestGenerateCronIntervals_SparseCombo(t *testing.T) {
	assert.Equal(
		t,
		[]map[string]int{
			{"Hour": 1, "Day": 2, "Month": 5},
			{"Hour": 1, "Day": 2, "Month": 9},
			{"Hour": 3, "Day": 2, "Month": 5},
			{"Hour": 3, "Day": 2, "Month": 9},
		},
		GenerateCronIntervals("", "1,3", "2", "5,9", ""),
	)
}

func TestGenerateCronIntervals_ItsAllStars(t *testing.T) {
	assert.Len(t, GenerateCronIntervals("*/2", "*/4", "*", "*/3", "*"), 30*6*4)
}
