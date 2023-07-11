package utils

import "regexp"

type xRegex struct {
}

var Regex = &xRegex{}

func (x *xRegex) Match(pattern string, s string) bool {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}
	return matched
}
