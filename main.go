package main

import (
	"flag"
	"fmt"
	"os"
)

const COMMAND_FIND_REPO = "find-repo"
const COMMAND_BUILD_REPO_INDEX = "build-repo-index"

func main() {
	executableName := os.Args[0]

	if len(os.Args) == 1 {
		printToolHelp(executableName)
	}

	command := os.Args[1]

	switch command {
	case "help":
		fallthrough
	case "-h":
		fallthrough
	case "--help":
		printToolHelp(executableName)
	case COMMAND_FIND_REPO:
		runFindRepoCommand(executableName)
	case COMMAND_BUILD_REPO_INDEX:
		runBuildRepoIndexCommand(executableName)
	default:
		printToolHelp(executableName)
	}
}

func printToolHelp(executableName string) {
	fmt.Println("Usage:")
	fmt.Printf("  %s command [options]\n", executableName)
	fmt.Println("Commands:")
	printSingleCommandDescription(COMMAND_FIND_REPO, "Finds repositories based on given filters")
	printSingleCommandDescription(COMMAND_BUILD_REPO_INDEX, "Builds the repository index")
	os.Exit(255)
}

func printSingleCommandDescription(command, commandHelp string) {
	fmt.Printf("\t%-18s %s\n", command, commandHelp)
}

func printCommandHelp(executableName string, command string, hasFilters bool) {
	filterText := func() string {
		if hasFilters {
			return " [<Filter Values>]"
		} else {
			return ""
		}
	}()
	fmt.Println("Usage:")
	fmt.Printf("  %s %s [options]%s\n", executableName, command, filterText)
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func runFindRepoCommand(executableName string) {
	var alfredOutput bool
	var showFindRepoHelp bool
	flag.BoolVar(&alfredOutput, "alfred", false, "Format output for Alfred")
	flag.BoolVar(&showFindRepoHelp, "h", false, "Show Help")
	flag.BoolVar(&showFindRepoHelp, "help", false, "Show Help")

	flag.CommandLine.Parse(os.Args[2:])

	if showFindRepoHelp {
		printCommandHelp(executableName, COMMAND_FIND_REPO, true)
	} else {
		positionalArgs := flag.Args()

		findRepo(executableName, alfredOutput, positionalArgs)
	}
}

func runBuildRepoIndexCommand(executableName string) {
	panic("unimplemented")
}
