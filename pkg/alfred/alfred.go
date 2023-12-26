package alfred

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	repo "ronkitay.com/griffin/pkg/repoindex"
)

func AsAlfred(matchingRepos []repo.RepoData) string {
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

func buildAlfredItem(repo repo.RepoData) Item {
	repoFullPath := filepath.Join(repo.BaseDir, repo.FullName)
	switch repo.Type {
	case "dir":
		return buildDirectoryLocation(repoFullPath, repo.FullName)
	case "archive":
		return buildArchiveLocation(repoFullPath, repo.FullName, repo.Url)
	case "gitlab":
		fallthrough
	case "github":
		return buildGitRepoLocation(repoFullPath, repo.FullName, repo.Url, repo.Type)
	default:
		panic("Unsupported locationType: " + repo.Type)
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
