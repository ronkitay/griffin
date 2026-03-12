package repoindex

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	config "ronkitay.com/griffin/pkg/configuration"
	csvHelper "ronkitay.com/griffin/pkg/csv"
)

type RepoData struct {
	BaseDir  string
	FullName string
	Url      string
	Type     string
}

func (datum RepoData) AsCsvRecord() []string {
	return []string{datum.BaseDir, datum.FullName, datum.Url, datum.Type}
}

func (datum RepoData) ToString() string {
	return filepath.Join(datum.BaseDir, datum.FullName)
}

func (datum RepoData) Matchable() string {
	return datum.FullName
}

func converter(noArchives bool, noDirs bool) func(csvData []string) (RepoData, error) {
	return func(csvData []string) (RepoData, error) {
		repoName := csvData[1]
		locationType := csvData[3]
		url := csvData[2]
		parentDir := csvData[0]

		switch locationType {
		case "dir":
			if !noDirs {
				return RepoData{BaseDir: parentDir, FullName: repoName, Type: "dir"}, nil
			}
		case "archive":
			if !noArchives {
				return RepoData{BaseDir: parentDir, FullName: repoName, Url: url, Type: locationType}, nil
			}
		case "gitlab":
			fallthrough
		case "github":
			return RepoData{BaseDir: parentDir, FullName: repoName, Url: url, Type: locationType}, nil
		}

		return RepoData{}, errors.New("Path skipped or not supported")
	}

}

func LoadIndex(noArchives bool, noDirs bool) []RepoData {
	return csvHelper.LoadIndex[RepoData](config.LoadConfiguration().RepoListLocation, converter(noArchives, noDirs))
}

func BuildRepoIndex() error {
	configuration := config.LoadConfiguration()
	configManager, err := config.NewConfigurationManager()
	if err != nil {
		return fmt.Errorf("error initializing configuration: %v", err)
	}

	roots, err := configManager.GetRepoRoots()
	if err != nil {
		return fmt.Errorf("error getting repository roots: %v", err)
	}

	processedRemotes := make(map[string]struct{})
	var repos []RepoData
	for _, rootLocation := range roots {
		reposFromRoot := locateRepos(rootLocation, processedRemotes)
		repos = append(repos, reposFromRoot...)
	}

	if err := csvHelper.SaveIndex(configuration.RepoListLocation, repos); err != nil {
		return fmt.Errorf("error saving repo index: %v", err)
	}

	return nil
}

func locateRepos(rootLocation string, processedRemotes map[string]struct{}) []RepoData {
	var repos []RepoData

	err := filepath.Walk(rootLocation, visit(rootLocation, &repos, processedRemotes))
	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", rootLocation, err)
	}

	repos = deDuplicate(repos)

	return repos
}

func visit(rootLocation string, paths *[]RepoData, processedRemotes map[string]struct{}) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}

		if info.IsDir() {
			// Check if the directory contains a .git directory
			gitPath := filepath.Join(path, ".git")
			_, err := os.Stat(gitPath)
			if err == nil {
				repoDir, repoName := dirAndName(rootLocation, path)
				remoteURL, err := getGitRemote(path)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					gitHttpUrl := gitURLToHTTP(remoteURL)
					repoType := repoType(gitHttpUrl)
					repoData := RepoData{BaseDir: repoDir, FullName: repoName, Url: gitHttpUrl, Type: repoType}
					*paths = append(*paths, repoData)

					if _, ok := processedRemotes[remoteURL]; !ok {
						processedRemotes[remoteURL] = struct{}{}
						worktrees, err := getWorktrees(path)
						if err == nil {
							for _, wtPath := range worktrees {
								if wtPath == path {
									continue
								}

								// Check if inside rootLocation
								rel, err := filepath.Rel(rootLocation, wtPath)
								if err == nil && !strings.HasPrefix(rel, "..") {
									// Inside
									wtDir, wtName := dirAndName(rootLocation, wtPath)
									*paths = append(*paths, RepoData{BaseDir: wtDir, FullName: wtName, Url: gitHttpUrl, Type: repoType})
								} else {
									// Outside - use repo name for the worktree
									*paths = append(*paths, RepoData{BaseDir: filepath.Dir(wtPath), FullName: repoName, Url: gitHttpUrl, Type: repoType})
								}
							}
						}
					}

					*paths = addParents(*paths, rootLocation, path)
					return filepath.SkipDir
				}
			}
		} else {
			if strings.HasSuffix(path, ".git") {
				file, err := os.Open(path)
				defer file.Close()
				if err == nil {

					scanner := bufio.NewScanner(file)

					if scanner.Scan() {
						firstLine := scanner.Text()
						re := regexp.MustCompile(`\s+`)
						cleanedLine := re.ReplaceAllString(firstLine, ";")
						archiveData := strings.Split(cleanedLine, ";")
						gitHttpUrl := gitURLToHTTP(archiveData[1])
						archiveDir, archiveName := dirAndName(rootLocation, path)
						repoData := RepoData{BaseDir: archiveDir, FullName: archiveName, Url: gitHttpUrl, Type: "archive"}
						*paths = append(*paths, repoData)
					}
				}
			}
		}

		return nil
	}
}

func deDuplicate(input []RepoData) []RepoData {
	encountered := map[RepoData]bool{}
	result := []RepoData{}

	for _, v := range input {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

func addParents(repos []RepoData, rootLocation, path string) []RepoData {
	if rootLocation == path {
		return repos
	} else {
		parentDir := filepath.Dir(path)
		dir, name := dirAndName(rootLocation, parentDir)
		repoData := RepoData{BaseDir: dir, FullName: name, Url: "-", Type: "dir"}
		repos = append(repos, repoData)
		return addParents(repos, rootLocation, parentDir)
	}
}

func dirAndName(rootLocation string, path string) (string, string) {
	if path == rootLocation {
		return filepath.Dir(path), filepath.Base(path)
	} else {
		return rootLocation, strings.Replace(path, rootLocation+"/", "", -1)
	}
}

func getGitRemote(dir string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get Git remote: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

var GIT_URL_REGEX = regexp.MustCompile(`git@([a-zA-Z0-9.-]+):`)

func gitURLToHTTP(url string) string {
	url = GIT_URL_REGEX.ReplaceAllString(url, "https://${1}/")
	url = strings.TrimSuffix(url, ".git")
	return url
}

func getWorktrees(dir string) ([]string, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %v", err)
	}

	var worktrees []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "worktree ") {
			path := strings.TrimPrefix(line, "worktree ")
			worktrees = append(worktrees, path)
		}
	}
	return worktrees, nil
}

func repoType(url string) string {
	if strings.Contains(url, "github") {
		return "github"
	}
	if strings.Contains(url, "gitlab") {
		return "gitlab"
	}
	return "unknown"
}
