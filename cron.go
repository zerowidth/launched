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

// Limit the cartesian product of cron fields
const MaxCronIntervals = 1440

func ValidateCronExpression(input string, min int, max int) bool {
	for part := range strings.SplitSeq(input, ",") {
		switch {
		case rangeRegex.MatchString(part):
			matches := rangeRegex.FindStringSubmatch(part)
			if !validNumber(matches[1], min, max) || !validNumber(matches[2], min, max) {
				return false
			}
		case starRegex.MatchString(part):
			matches := starRegex.FindStringSubmatch(part)
			if matches[2] != "" && !validDivisor(matches[2]) {
				return false
			}
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

func validDivisor(input string) bool {
	divisor, err := strconv.Atoi(input)
	return err == nil && divisor > 0
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
	choices := cronIntervalChoices(minute, hour, day_of_month, month, weekday)

	// Validation rejects anything over the limit on create, but a plist stored
	// before the limit existed would otherwise allocate its way to an OOM every
	// time it's rendered.
	if countCronIntervals(choices) > MaxCronIntervals {
		return nil
	}

	intervals := combineIntervals(choices)

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

func cronIntervalChoices(minute, hour, day_of_month, month, weekday string) []intervalChoices {
	return []intervalChoices{
		{"Minute", cronIntervals(minute, 0, 59)},
		{"Hour", cronIntervals(hour, 0, 23)},
		{"Day", cronIntervals(day_of_month, 1, 31)},
		{"Month", cronIntervals(month, 1, 12)},
		{"Weekday", cronIntervals(weekday, 0, 6)},
	}
}

// CronIntervalCount reports how many calendar intervals these expressions would
// generate without building them.
func CronIntervalCount(minute, hour, day_of_month, month, weekday string) int {
	return countCronIntervals(cronIntervalChoices(minute, hour, day_of_month, month, weekday))
}

func countCronIntervals(choices []intervalChoices) int {
	count := 0
	for _, choice := range choices {
		if len(choice.intervals) == 0 {
			continue
		}
		if count == 0 {
			count = len(choice.intervals)
		} else {
			count *= len(choice.intervals)
		}
	}
	return count
}

func cronIntervals(cron string, min, max int) []int {
	set := map[int]struct{}{}

	for part := range strings.SplitSeq(cron, ",") {
		switch {
		case rangeRegex.MatchString(part):
			matches := rangeRegex.FindStringSubmatch(part)
			start, _ := strconv.Atoi(matches[1])
			end, _ := strconv.Atoi(matches[2])
			// clamp to min/max to prevent invalid intervals
			if start < min {
				start = min
			}
			if end > max {
				end = max
			}
			for i := start; i <= end; i++ {
				set[i] = struct{}{}
			}
		case starRegex.MatchString(part):
			matches := starRegex.FindStringSubmatch(part)
			divisor := 1
			if matches[2] != "" {
				divisor, _ = strconv.Atoi(matches[2])
			}
			// `*` and `*/1` are unconstrained; `*/0` would divide by zero
			if divisor < 2 {
				continue
			}
			for i := min; i <= max; i++ {
				if i%divisor == 0 {
					set[i] = struct{}{}
				}
			}
		case numberRegex.MatchString(part):
			number, err := strconv.Atoi(part)
			if err != nil || number < min || number > max {
				continue
			}
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
