package finder

import (
	"fmt"
	"path/filepath"

	alfred "ronkitay.com/griffin/pkg/alfred"
	matcher "ronkitay.com/griffin/pkg/matcher"
	projectIndex "ronkitay.com/griffin/pkg/projectindex"
	repo "ronkitay.com/griffin/pkg/repoindex"
)

func FindRepo(executableName string, noArchives bool, noDirs bool, alfredOutput bool, args []string) {
	allRepos := repo.LoadIndex(noArchives, noDirs)

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

func FindProjects(executableName string, alfredOutput bool, args []string) {
	allProjects := projectIndex.LoadIndex()

	regexPattern := matcher.BuildPattern(args)

	matchingProjects := matcher.MatchProjects(allProjects, regexPattern)

	printProjectPaths(matchingProjects)

}

func printProjectPaths(matchingProjects []projectIndex.ProjectData) {
	for _, project := range matchingProjects {
		fmt.Println(filepath.Join(project.BaseDir, project.FullName))
	}
}
