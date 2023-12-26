# GRIFFIN

GRIFFIN - **G**it **R**epository **I**ndexer for **F**uzzy **F**inding and **IN**spection

A CLI tool for indexing all git repos and code projects on your computer

## Usage

### Configuring

Create a configuration file at `~/.config/griffin/config.json`

Configure the paths to be indexed.

Example:

```json
{
	"repoRoots": [
		"${HOME}/personal",
		"${HOME}/work"
	]
}
```

### Building a Repository Index

```bash
griffin build-repo-index
```

### Searching for Repos

```bash
griffin find-repo [-alfred] [search arguments]
```

### Shell Integration

Add the following to your ~/.zshrc

```bash
function r() {
    TEMP_LIST_FILE=$(mktemp)

    ${HOME}/tools/griffin find-repo $* > ${TEMP_LIST_FILE}

    if [[ "$(cat ${TEMP_LIST_FILE} | wc -l)" -eq "1" ]]; 
    then 
        DIR_TO_SWITCH_TO=`cat ${TEMP_LIST_FILE}`
    else
        DIR_TO_SWITCH_TO=`cat ${TEMP_LIST_FILE} | fzf --preview 'tree -L 2 -C {}'`
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
```

**Note:** Make sure you have `tree` and `fzf` installed for this command to work.

```bash
brew install fzf
brew install tree
```
