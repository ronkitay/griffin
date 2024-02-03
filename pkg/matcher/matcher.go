package matcher

import (
	"os"
	"regexp"
	"strings"
)

func BuildPattern(args []string) *regexp.Regexp {
	var searchPattern string

	if len(args) == 0 {
		searchPattern = ".*"
	} else {
		searchPattern = "(?i).*" + strings.Join(args, ".*") + ".*"
	}

	regexPattern, regexError := regexp.Compile(searchPattern)
	if regexError != nil {
		os.Exit(3)
	}

	return regexPattern
}

type Matchable interface {
	Matchable() string
}

func MatchItems[T Matchable](elements []T, regexPattern *regexp.Regexp) []T {
	var result []T

	for _, element := range elements {
		if regexPattern.MatchString(element.Matchable()) {
			result = append(result, element)
		}
	}

	return result
}
