package everything

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func asAlfred(matchingRepos []RepoData) string {
	var items []Item

	for _, repo := range matchingRepos {
		items = append(items, buildAlfredItem(repo))
	}

	result := map[string][]Item{
		"items": items,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		os.Exit(1)
	}

	return string(jsonData)
}

func buildAlfredItem(repo RepoData) Item {
	repoFullPath := filepath.Join(repo.repoDir, repo.repoName)
	switch repo.locationType {
	case "dir":
		return buildDirectoryLocation(repoFullPath, repo.repoName)
	case "archive":
		return buildArchiveLocation(repoFullPath, repo.repoName, repo.url)
	case "gitlab":
		fallthrough
	case "github":
		return buildGitRepoLocation(repoFullPath, repo.repoName, repo.url, repo.locationType)
	default:
		panic("Unsupported locationType: " + repo.locationType)
	}
}

func buildDirectoryLocation(repoDir string, repoName string) Item {
	return Item{
		Valid:    true,
		UID:      repoName,
		Title:    repoName,
		Subtitle: "Open in TERMINAL (🖥️) : " + repoDir,
		Arg:      repoDir,
		Icon: Icon{
			Path: "icons/dir.jpg",
		},
	}
}

func buildArchiveLocation(repoDir string, repoName string, url string) Item {
	return Item{
		Valid:    true,
		UID:      repoName,
		Title:    repoName,
		Subtitle: "Open in TERMINAL (🖥️) : " + repoDir,
		Arg:      repoDir,
		Mods: map[string]Modifier{
			"alt": {
				Valid:    true,
				Arg:      url,
				Subtitle: "Open in WEB (☁️): " + url,
			},
		},
		Icon: Icon{
			Path: "icons/archive.jpg",
		},
	}
}

func buildGitRepoLocation(repoDir string, repoName string, url string, locationType string) Item {
	return Item{
		Valid:    true,
		UID:      repoName,
		Title:    repoName,
		Subtitle: "Open in TERMINAL (🖥️) : " + repoDir,
		Arg:      repoDir,
		Mods: map[string]Modifier{
			"alt": {
				Valid:    true,
				Arg:      url,
				Subtitle: "Open in WEB (☁️): " + url,
			},
			"ctrl": {
				Valid:    true,
				Arg:      repoDir,
				Subtitle: "Open in EDITOR (📝): " + repoDir,
			},
		},
		Icon: Icon{
			Path: "icons/" + locationType + ".jpg",
		},
	}
}

type Item struct {
	Valid    bool                `json:"valid"`
	UID      string              `json:"uid"`
	Title    string              `json:"title"`
	Subtitle string              `json:"subtitle"`
	Arg      string              `json:"arg"`
	Mods     map[string]Modifier `json:"mods"`
	Icon     Icon                `json:"icon"`
}

type Modifier struct {
	Valid    bool   `json:"valid"`
	Arg      string `json:"arg"`
	Subtitle string `json:"subtitle"`
}

type Icon struct {
	Path string `json:"path"`
}
