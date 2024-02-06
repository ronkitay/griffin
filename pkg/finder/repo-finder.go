package finder

import (
	"fmt"

	alfred "ronkitay.com/griffin/pkg/alfred"
	matcher "ronkitay.com/griffin/pkg/matcher"
	projectIndex "ronkitay.com/griffin/pkg/projectindex"
	repo "ronkitay.com/griffin/pkg/repoindex"
)

func FindRepo(executableName string, noArchives bool, noDirs bool, alfredOutput bool, args []string) {
	allRepos := repo.LoadIndex(noArchives, noDirs)

	regexPattern := matcher.BuildPattern(args)

	matchingRepos := matcher.MatchItems(allRepos, regexPattern)

	if alfredOutput {
		result := alfred.ReposAsAlfred(matchingRepos)
		fmt.Println(result)
	} else {
		printPaths(matchingRepos)
	}
}

type Printable interface {
	ToString() string
}

func printPaths[T Printable](matchingItems []T) {
	for _, item := range matchingItems {
		fmt.Println(item.ToString())
	}
}

func FindProjects(executableName string, alfredOutput bool, args []string) {
	allProjects := projectIndex.LoadIndex()

	regexPattern := matcher.BuildPattern(args)

	matchingProjects := matcher.MatchItems(allProjects, regexPattern)

	if alfredOutput {
		result := alfred.ProjectsAsAlfred(matchingProjects)
		fmt.Println(result)
	} else {
		printPaths(matchingProjects)
	}

}
