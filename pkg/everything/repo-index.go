package everything

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	config "ronkitay.com/griffin/pkg/configuration"
)

type RepoData struct {
	repoDir      string
	repoName     string
	url          string
	locationType string
}

func loadIndex() []RepoData {
	csvData := loadIndexCsv()

	var items []RepoData

	for _, line := range csvData {
		repoName := line[1]
		locationType := line[3]
		url := line[2]
		parentDir := line[0]

		switch locationType {
		case "dir":
			items = append(items, RepoData{repoDir: parentDir, repoName: repoName, locationType: "dir"})
		case "archive":
			fallthrough
		case "gitlab":
			fallthrough
		case "github":
			items = append(items, RepoData{repoDir: parentDir, repoName: repoName, url: url, locationType: locationType})
		}
	}
	return items
}

func loadIndexCsv() [][]string {
	configuration := config.LoadConfiguration()

	file, fileOpenError := os.Open(configuration.RepoListLocation)
	if fileOpenError != nil {
		os.Exit(1)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	data, csvReadError := csvReader.ReadAll()
	if csvReadError != nil {
		os.Exit(2)
	}
	return data
}

func buildRepoIndex() {
	userHomeDir, _ := os.UserHomeDir()
	configuration := config.LoadConfiguration()
	// fmt.Println(configuration)

	var repos []RepoData
	for _, rootLocation := range configuration.UserConfiguration.RepoRoots {
		interpolatedRootLocation := strings.Replace(rootLocation, "${HOME}", userHomeDir, -1)
		reposFromRoot := locateRepos(interpolatedRootLocation)
		repos = append(repos, reposFromRoot...)
	}

	file, err := os.Create(configuration.RepoListLocation + ".new")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'

	// Write data
	for _, repo := range repos {
		row := []string{repo.repoDir, repo.repoName, repo.url, repo.locationType}
		err := writer.Write(row)
		if err != nil {
			fmt.Println("Error writing CSV row:", err)
			return
		}
	}

	writer.Flush()

	// Check for errors during Flush
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing CSV writer:", err)
		return
	}

	os.Rename(configuration.RepoListLocation+".new", configuration.RepoListLocation)
}

func locateRepos(rootLocation string) []RepoData {
	var repos []RepoData

	err := filepath.Walk(rootLocation, visit(rootLocation, &repos))
	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", rootLocation, err)
	}

	repos = deDuplicate(repos)
	// fmt.Println("List of paths:")
	// for _, repo := range repos {
	// 	fmt.Println(repo)
	// }

	return repos
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

func visit(rootLocation string, paths *[]RepoData) filepath.WalkFunc {
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
					repoData := RepoData{repoDir: repoDir, repoName: repoName, url: gitHttpUrl, locationType: repoType}
					*paths = append(*paths, repoData)

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
						repoData := RepoData{repoDir: archiveDir, repoName: archiveName, url: gitHttpUrl, locationType: "archive"}
						*paths = append(*paths, repoData)
					}
				}
			}
		}

		return nil
	}
}

func addParents(repos []RepoData, rootLocation, path string) []RepoData {
	if rootLocation == path {
		return repos
	} else {
		parentDir := filepath.Dir(path)
		dir, name := dirAndName(rootLocation, parentDir)
		repoData := RepoData{repoDir: dir, repoName: name, url: "-", locationType: "dir"}
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

func repoType(url string) string {
	if strings.Contains(url, "github") {
		return "github"
	}
	if strings.Contains(url, "gitlab") {
		return "gitlab"
	}
	return "unknown"
}
