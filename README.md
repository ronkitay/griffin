# GRIFFIN

GRIFFIN - **G**it **R**epository **I**ndexer for **F**uzzy **F**inding and **IN**spection

A CLI tool for indexing all git repos and code projects on your computer

## Installation

### Quick Install (One-liner)

```bash
curl -sSL https://raw.githubusercontent.com/ronkitay/griffin/main/install.sh | sh
```

This will:
- Download the latest griffin binary for your system
- Install it to `${HOME}/tools/`
- Add `${HOME}/tools` to your PATH in `~/.zshrc` (if not already present)
- Remove the quarantine attribute on macOS

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
    },
    "disableHttpProbe": false
}
```

### Building a Repository Index

```bash
griffin build-repo-index [--no-http-probe]
```

#### Self-Hosted Git Repository Detection

When indexing repositories whose remote host domain does not contain `github` or `gitlab`, Griffin determines the repository type in the following order:

1. **HTTP Host Probe (cached once per domain):** Sends an unauthenticated `HEAD` request (with a 2-second timeout and no redirects) to the remote host origin and checks for GitHub (`X-GitHub-Request-Id` / `Server: github.com`) or GitLab (`X-Gitlab-Meta`) headers.
2. **Local Filesystem Fallback:** If the host is unreachable (or HTTP probing is disabled), checks the repository root for `.github` (`github`) or `.gitlab-ci.yml` / `.gitlab` (`gitlab`).
3. **Generic Git Fallback:** Defaults to `git` if neither matches.

To disable outbound HTTP probing during indexing:
- Pass `--no-http-probe` for a single run: `griffin build-repo-index --no-http-probe`
- Or persist it in `~/.config/griffin/config.json`: `griffin configure --disable-http-probe=true`

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
