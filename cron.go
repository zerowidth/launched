package main

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var rangeRegex = regexp.MustCompile(`^(\d+)-(\d+)$`)
var starRegex = regexp.MustCompile(`^\*(/(\d+))?$`)
var numberRegex = regexp.MustCompile(`^\d+$`)

func ValidateCronExpression(input string, min int, max int) bool {
	parts := strings.Split(input, ",")
	for _, part := range parts {
		switch {
		case rangeRegex.MatchString(part):
			matches := rangeRegex.FindStringSubmatch(part)
			if !validNumber(matches[1], min, max) || !validNumber(matches[2], min, max) {
				return false
			}
		case starRegex.MatchString(part):
			// fine i guess, if divisor is greater than max it'll still hit interval 0
		case numberRegex.MatchString(part):
			if !validNumber(part, min, max) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func validNumber(input string, min, max int) bool {
	number, err := strconv.Atoi(input)
	if err != nil {
		return false
	}
	if number < min || number > max {
		return false
	}
	return true
}

type intervalChoices struct {
	name      string
	intervals []int
}

type choice struct {
	name     string
	interval int
}

// assume the input is valid
func GenerateCronIntervals(minute, hour, day_of_month, month, weekday string) []map[string]int {
	intervals := combineIntervals([]intervalChoices{
		{"Minute", cronIntervals(minute, 0, 59)},
		{"Hour", cronIntervals(hour, 0, 23)},
		{"DayOfMonth", cronIntervals(day_of_month, 1, 31)},
		{"Month", cronIntervals(month, 1, 12)},
		{"Weekday", cronIntervals(weekday, 0, 6)},
	})

	sum := []map[string]int{}
	for _, interval := range intervals {
		combination := map[string]int{}
		for _, choice := range interval {
			combination[choice.name] = choice.interval
		}
		sum = append(sum, combination)
	}

	return sum
}

func cronIntervals(cron string, min, max int) []int {
	set := map[int]struct{}{}

	parts := strings.Split(cron, ",")
	for _, part := range parts {
		switch {
		case rangeRegex.MatchString(part):
			matches := rangeRegex.FindStringSubmatch(part)
			start, _ := strconv.Atoi(matches[1])
			end, _ := strconv.Atoi(matches[2])
			for i := start; i <= end; i++ {
				set[i] = struct{}{}
			}
		case starRegex.MatchString(part):
			matches := starRegex.FindStringSubmatch(part)
			divisor := 1
			if matches[2] != "" {
				divisor, _ = strconv.Atoi(matches[2])
			}
			for i := min; i <= max; i++ {
				if i%divisor == 0 {
					set[i] = struct{}{}
				}
			}
		case numberRegex.MatchString(part):
			number, _ := strconv.Atoi(part)
			set[number] = struct{}{}
		}
	}
	values := []int{}
	for v := range set {
		values = append(values, v)
	}
	sort.Ints(values)
	return values
}

func combineIntervals(choices []intervalChoices) [][]choice {
	return combine(choices, [][]choice{})
}

func combine(choices []intervalChoices, chosen [][]choice) [][]choice {
	if len(choices) == 0 {
		return chosen
	}
	if len(choices[0].intervals) == 0 {
		return combine(choices[1:], chosen)
	}
	combined := [][]choice{}
	name := choices[0].name
	for _, interval := range choices[0].intervals {
		if len(chosen) == 0 {
			combined = append(combined, combine(choices[1:], [][]choice{{{name, interval}}})...)
		} else {
			for _, cs := range chosen {
				combined = append(combined, combine(choices[1:], [][]choice{append(cs, choice{name, interval})})...)
			}
		}
	}
	return combined
}
