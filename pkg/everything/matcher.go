package everything

import (
	"os"
	"regexp"
	"strings"
)

func buildPattern(args []string) *regexp.Regexp {
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

func matchRepos(repoList []RepoData, regexPattern *regexp.Regexp) []RepoData {
	var result []RepoData

	for _, element := range repoList {
		if regexPattern.MatchString(element.repoName) {
			result = append(result, element)
		}
	}

	return result
}
