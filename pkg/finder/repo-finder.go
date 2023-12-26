package finder

import (
	"fmt"
	"path/filepath"
	alfred "ronkitay.com/griffin/pkg/alfred"
	repo "ronkitay.com/griffin/pkg/repoindex"
	matcher "ronkitay.com/griffin/pkg/matcher"
)

func FindRepo(executableName string, alfredOutput bool, args []string) {
	allRepos := repo.LoadIndex()

	regexPattern := matcher.BuildPattern(args)

	matchingRepos := matcher.MatchRepos(allRepos, regexPattern)

	if alfredOutput {
		result := alfred.AsAlfred(matchingRepos)
		fmt.Println(result)
	} else {
		printPaths(matchingRepos)
	}
}

func printPaths(matchingRepos []repo.RepoData) {
	for _, repo := range matchingRepos {
		fmt.Println(filepath.Join(repo.BaseDir, repo.FullName))
	}
}
