package shell

import (
	"fmt"
	"os"
	"os/exec"

	"ronkitay.com/griffin/pkg/terminal"
)

func GenerateIntegration() {

	dependencies := []string{"fzf", "tree"}

	for _, tool := range dependencies {
		_, notInstalledError := toolIsInstalled(tool)
		if notInstalledError != nil {
			fmt.Printf("Tool '%s' not found in the PATH.\nInstall it using the following command:\nbrew install %s\n\n", tool, tool)
			os.Exit(1)
		}
	}

	scriptTemplate := `
	function r() {
		TEMP_LIST_FILE=$(mktemp)
	
		griffin find-repo $* > "${TEMP_LIST_FILE}"
	
		if [[ "$(cat "${TEMP_LIST_FILE}" | wc -l)" -eq "1" ]]; 
		then 
			DIR_TO_SWITCH_TO=$(cat "${TEMP_LIST_FILE}")
		else
			DIR_TO_SWITCH_TO=$(cat "${TEMP_LIST_FILE}" | fzf +m --preview 'tree -L 2 -C {}')
		fi
		rm "${TEMP_LIST_FILE}"
	
		if [[ "${DIR_TO_SWITCH_TO}" = *.git ]]; 
		then
			cd $(dirname "${DIR_TO_SWITCH_TO}");
			echo "%s%sTo access the repo, run the following command:%s"
			echo ""
			echo "unarchive $(basename ${DIR_TO_SWITCH_TO})"
			echo ""
		else
			cd "${DIR_TO_SWITCH_TO}"
		fi
	}
	
	function or() {
		TEMP_LIST_FILE=$(mktemp)
		
		griffin find-repo --noarchive --nodir $* > "${TEMP_LIST_FILE}"
		
		PROJECT_DIR=$(cat "${TEMP_LIST_FILE}" | fzf +m --preview 'tree -L 2 -C {}')
		
		rm "${TEMP_LIST_FILE}"
		
		if [[ -n "${PROJECT_DIR}" ]]; then
			griffin open-in-ide "${PROJECT_DIR}"
		fi
	}

	function p() {
		TEMP_LIST_FILE=$(mktemp)

		griffin find-project $* > "${TEMP_LIST_FILE}"
	
		if [[ "$(cat "${TEMP_LIST_FILE}" | wc -l)" -eq "1" ]]; 
		then 
			DIR_TO_SWITCH_TO=$(cat "${TEMP_LIST_FILE}")
		else
			DIR_TO_SWITCH_TO=$(cat "${TEMP_LIST_FILE}" | fzf +m --preview 'tree -L 2 -C {}')
		fi
		rm "${TEMP_LIST_FILE}"
	
		cd "${DIR_TO_SWITCH_TO}"
	}
	`

	script := fmt.Sprintf(scriptTemplate, terminal.BOLD_COLOR, terminal.GREEN_COLOR, terminal.RESET_COLORS)

	fmt.Println(script)
}

func toolIsInstalled(tool string) (bool, error) {
	if _, toolNotInPathError := exec.LookPath(tool); toolNotInPathError != nil {
		return false, toolNotInPathError
	} else {
		return true, nil
	}
}
