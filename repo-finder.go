package main

import (
	"fmt"
	"path/filepath"
)

func findRepo(executableName string, alfredOutput bool, args []string) {
	allRepos := loadIndex()

	regexPattern := buildPattern(args)

	matchingRepos := matchRepos(allRepos, regexPattern)

	if alfredOutput {
		result := asAlfred(matchingRepos)
		fmt.Println(result)
	} else {
		printPaths(matchingRepos)
	}
}

func printPaths(matchingRepos []RepoData) {
	for _, repo := range matchingRepos {
		fmt.Println(filepath.Join(repo.repoDir, repo.repoName))
	}
}
