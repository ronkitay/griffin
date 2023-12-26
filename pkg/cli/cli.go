package cli

import (
	"flag"
	"fmt"
	"os"

	"ronkitay.com/griffin/pkg/repoindex"
	"ronkitay.com/griffin/pkg/finder"
	"ronkitay.com/griffin/pkg/idelauncher"
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
	{"open-in-ide", "Opens a given path in the appropriate IDE", runInIDECommand},
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

	var noArchives bool
	flag.BoolVar(&noArchives, "noarchive", false, "Filter out Archives")

	var noDirs bool
	flag.BoolVar(&noDirs, "nodir", false, "Filter out Directories")

	flag.CommandLine.Parse(os.Args[2:])

	if showFindRepoHelp {
		printCommandHelp(executableName, command.name, true)
	} else {
		positionalArgs := flag.Args()

		finder.FindRepo(executableName, noArchives, noDirs, alfredOutput, positionalArgs)
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
	script := `
	function r() {
		TEMP_LIST_FILE=$(mktemp)
	
		${HOME}/tools/griffin find-repo $* > ${TEMP_LIST_FILE}
	
		if [[ "$(cat ${TEMP_LIST_FILE} | wc -l)" -eq "1" ]]; 
		then 
			DIR_TO_SWITCH_TO=$(cat ${TEMP_LIST_FILE})
		else
			DIR_TO_SWITCH_TO=$(cat ${TEMP_LIST_FILE} | fzf --preview 'tree -L 2 -C {}')
		fi
		rm ${TEMP_LIST_FILE}
	
		if [[ "${DIR_TO_SWITCH_TO}" = *.git ]]; 
		then
			cd $(dirname ${DIR_TO_SWITCH_TO});
			echo "${BRIGHT}${GREEN}To access the repo, run the following command:${NORMAL}"
			echo ""
			echo "unarchive $(basename ${DIR_TO_SWITCH_TO})"
			echo ""
		else
			cd $DIR_TO_SWITCH_TO
		fi
	}
	
	function o() {
		TEMP_LIST_FILE=$(mktemp)
		${HOME}/tools/griffin find-repo --noarchive --nodir $* > ${TEMP_LIST_FILE}
		PROJECT_DIR=$(cat ${TEMP_LIST_FILE} | fzf --preview 'tree -L 2 -C {}')
		rm ${TEMP_LIST_FILE}
		if [[ -n "${PROJECT_DIR}" ]]; then
			griffin open-in-ide $PROJECT_DIR
		fi
	}
	`

	fmt.Println(script)

}
func runConfigureCommand(command *Command, executableName string) {
	panic(command.name + " not implemented yet!")
}
func runInIDECommand(command *Command, executableName string) {
	idelauncher.OpenInIDE(os.Args[2])
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
