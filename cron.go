package main

import (
	"regexp"
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
