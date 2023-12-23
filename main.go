package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		printHelp(os.Args[0])
	}

	command := os.Args[1]

	switch command {
	case "help":
		printHelp(os.Args[0])
	case "find-repo":
		var alfredOutput bool
		var showFindRepoHelp bool
		flag.BoolVar(&alfredOutput, "alfred", false, "Format output for Alfred")
		flag.BoolVar(&showFindRepoHelp, "h", false, "Show Help")
		flag.BoolVar(&showFindRepoHelp, "help", false, "Show Help")

		flag.CommandLine.Parse(os.Args[2:])

		positionalArgs := flag.Args()

		findRepo(os.Args[0], alfredOutput, positionalArgs)
	default:
		fmt.Println("Got command" + command)
		return
	}
}

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
		fmt.Println(repo.repoDir)
	}
}

func printHelp(executableName string) {
	fmt.Printf("%s\n", executableName)
	fmt.Printf("Missing command!\n")
	os.Exit(255)
}
