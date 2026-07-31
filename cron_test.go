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

func TestValidateCronExpression_ZeroDivisor(t *testing.T) {
	// `*/0` used to validate, then divide by zero when the plist was rendered
	assert.False(t, ValidateCronExpression("*/0", 0, 59))
	assert.True(t, ValidateCronExpression("*/1", 0, 59))
	assert.Empty(t, GenerateCronIntervals("*/1", "", "", "", ""))
}

func TestCronIntervals_ClampsToBounds(t *testing.T) {
	// out-of-range input shouldn't be able to allocate its way out of the range
	assert.Equal(t, []int{58, 59}, cronIntervals("58-99999999", 0, 59))
	assert.Empty(t, cronIntervals("100-200", 0, 59))
	assert.Empty(t, cronIntervals("99", 0, 59))
}

func TestCronIntervalCount(t *testing.T) {
	assert.Equal(t, 0, CronIntervalCount("", "", "", "", ""))
	assert.Equal(t, 0, CronIntervalCount("*", "*", "*", "*", "*"))
	assert.Equal(t, 1, CronIntervalCount("30", "4", "", "", ""))
	assert.Equal(t, 60, CronIntervalCount("0-59", "", "", "", ""))
	assert.Equal(t, 1440, CronIntervalCount("0-59", "0-23", "", "", ""))
	assert.Equal(t, 3749760, CronIntervalCount("0-59", "0-23", "1-31", "1-12", "0-6"))
}

func TestGenerateCronIntervals_RefusesToBlowUp(t *testing.T) {
	// the cartesian product of every field is 3.7M entries and gigabytes of XML
	assert.Nil(t, GenerateCronIntervals("0-59", "0-23", "1-31", "1-12", "0-6"))
	assert.Len(t, GenerateCronIntervals("0-59", "0-23", "", "", ""), MaxCronIntervals)
}

func TestValidateRejectsTooManyIntervals(t *testing.T) {
	broad := LaunchdPlist{Name: "x", Command: "x", Minute: "0-59", Hour: "0-23", DayOfMonth: "1-31", Month: "1-12", Weekday: "0-6"}
	errors := broad.Validate()
	assert.Contains(t, errors, "LaunchdPlist.Minute")

	ok := LaunchdPlist{Name: "x", Command: "x", Minute: "0-59", Hour: "0-23"}
	assert.Nil(t, ok.Validate())
}
