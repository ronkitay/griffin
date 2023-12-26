package matcher

import (
	"os"
	"regexp"
	"strings"

	repo "ronkitay.com/griffin/pkg/repoindex"
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

func MatchRepos(repoList []repo.RepoData, regexPattern *regexp.Regexp) []repo.RepoData {
	var result []repo.RepoData

	for _, element := range repoList {
		if regexPattern.MatchString(element.FullName) {
			result = append(result, element)
		}
	}

	return result
}
