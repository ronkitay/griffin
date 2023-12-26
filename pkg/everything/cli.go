package everything

import (
	"flag"
	"fmt"
	"os"

	"ronkitay.com/griffin/pkg/repoindex"
)

type CommandHandler func(*Command, string)

type Command struct {
	name, description string
	handler           CommandHandler
}

var COMMANDS = []Command{
	{"find-repo", "Finds repositories based on given filters", runFindRepoCommand},
	{"build-repo-index", "Builds the repository index", runBuildRepoIndexCommand},
	{"find-project", "Finds projects based on given filters", runFindProjectCommand},
	{"build-project-index", "Builds the projects index", runBuildProjectIndexCommand},
	{"shell-integration", "Generates Shell Integration commands", runShellIntegrationCommand},
	{"configure", "Configure the tool", runConfigureCommand},
}

const (
	RESET_COLORS = "\033[0m"
	BOLD_COLOR   = "\033[1m"

	RED_COLOR   = "\033[31m"
	GREEN_COLOR = "\033[32m"
	WHITE_COLOR = "\033[37m"
)

const COMMAND_NOT_SUPPORTED_ERROR_MESSAGE = BOLD_COLOR + RED_COLOR + "Command '" + WHITE_COLOR + "%s" + RED_COLOR + "' is not supported!" + RESET_COLORS + "\n"

func Run() {
	executableName := os.Args[0]

	if len(os.Args) == 1 {
		printToolHelp(executableName)
	}

	commandName := os.Args[1]

	if userRequestsHelp(commandName) {
		printToolHelp(executableName)
	} else {
		command := matchCommand(commandName)
		if command != nil {
			command.handler(command, executableName)
		} else {
			fmt.Fprintf(os.Stderr, COMMAND_NOT_SUPPORTED_ERROR_MESSAGE, commandName)
			printToolHelp(executableName)
		}
	}
}

func userRequestsHelp(commandName string) bool {
	return commandName == "-h" || commandName == "--help" || commandName == "help"
}

func printToolHelp(executableName string) {
	fmt.Println("Usage:")
	fmt.Printf("  %s commandName [options]\n", executableName)
	fmt.Println("Commands:")
	for _, commandName := range COMMANDS {
		printSingleCommandDescription(commandName.name, commandName.description)
	}
	os.Exit(255)
}

func printSingleCommandDescription(commandName, commandHelp string) {
	fmt.Printf("\t"+GREEN_COLOR+"%-20s"+RESET_COLORS+" %s\n", commandName, commandHelp)
}

func matchCommand(commandName string) *Command {
	for _, potentialMatch := range COMMANDS {
		if potentialMatch.name == commandName {
			return &potentialMatch
		}
	}
	return nil
}

func runFindRepoCommand(command *Command, executableName string) {
	var showFindRepoHelp bool
	flag.BoolVar(&showFindRepoHelp, "h", false, "Show Help")
	flag.BoolVar(&showFindRepoHelp, "help", false, "Show Help")

	var alfredOutput bool
	flag.BoolVar(&alfredOutput, "alfred", false, "Format output for Alfred")

	flag.CommandLine.Parse(os.Args[2:])

	if showFindRepoHelp {
		printCommandHelp(executableName, command.name, true)
	} else {
		positionalArgs := flag.Args()

		findRepo(executableName, alfredOutput, positionalArgs)
	}
}

func runBuildRepoIndexCommand(command *Command, executableName string) {
	repoindex.BuildRepoIndex()
}
func runFindProjectCommand(command *Command, executableName string) {
	panic(command.name + " not implemented yet!")
}
func runBuildProjectIndexCommand(command *Command, executableName string) {
	panic(command.name + " not implemented yet!")
}
func runShellIntegrationCommand(command *Command, executableName string) {
	panic(command.name + " not implemented yet!")
}
func runConfigureCommand(command *Command, executableName string) {
	panic(command.name + " not implemented yet!")
}

func printCommandHelp(executableName string, commandName string, hasFilters bool) {
	filterText := func() string {
		if hasFilters {
			return " [<Filter Values>]"
		} else {
			return ""
		}
	}()
	fmt.Println("Usage:")
	fmt.Printf("  %s %s [options]%s\n", executableName, commandName, filterText)
	fmt.Println("Options:")
	flag.PrintDefaults()
}
