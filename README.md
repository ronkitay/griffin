# GRIFFIN

GRIFFIN - **G**it **R**epository **I**ndexer for **F**uzzy **F**inding and **IN**spection

A CLI tool for indexing all git repos and code projects on your computer

## Usage

### Configuring

Create a configuration file at `~/.config/griffin/config.json`

Configure the paths to be indexed.
Configure the IDEs to be used per programming language. (Currently supports go, java, kotlin, python, node, and rust)

Example:

```json
{
    "repoRoots": [
        "${HOME}/personal",
        "${HOME}/work"
    ],
    "ideConfiguration": {
        "default": "Visual Studio Code.app", 
        "go": "GoLand.app",
        "java": "IntelliJ IDEA CE.app",
        "python": "PyCharm CE.app",
        "node": "WebStorm.app",
        "rust": "RustRover.app"
    }
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

### Opening a path in an IDE

```bash
griffin open-in-ide <path>
```

### Shell Integration

Add the following to your ~/.zshrc

```bash
source <(griffin shell-integration)
```

**Note:** Make sure you have `tree` and `fzf` installed for this command to work.

```bash
brew install fzf
brew install tree
```
