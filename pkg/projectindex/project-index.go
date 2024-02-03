package projectindex

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	config "ronkitay.com/griffin/pkg/configuration"
	csvHelper "ronkitay.com/griffin/pkg/csv"
	repoIndex "ronkitay.com/griffin/pkg/repoindex"
)

type ProjectData struct {
	BaseDir  string
	FullName string
	Type     string
}

func (datum ProjectData) AsCsvRecord() []string {
	return []string{datum.BaseDir, datum.FullName, datum.Type}
}

func FromCsvRecord(data []string) (ProjectData, error) {
	return ProjectData{
		BaseDir:  data[0],
		FullName: data[1],
		Type:     data[2],
	}, nil
}

func LoadIndex() []ProjectData {
	return csvHelper.LoadIndex[ProjectData](config.LoadConfiguration().ProjectListLocation, FromCsvRecord)
}

func BuildProjectIndex() {

	repos := repoIndex.LoadIndex(true, true)

	var projects []ProjectData

	for _, repo := range repos {
		repoRoot := filepath.Join(repo.BaseDir, repo.FullName)

		scanRepoForProjects(repoRoot, &projects)
	}

	csvHelper.SaveIndex(config.LoadConfiguration().ProjectListLocation, projects)
}

func scanRepoForProjects(rootLocation string, projects *[]ProjectData) {
	err := filepath.WalkDir(rootLocation, visitDirs(rootLocation, projects))

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", rootLocation, err)
	}
}

func visitDirs(rootLocation string, projects *[]ProjectData) fs.WalkDirFunc {
	return func(path string, info os.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}

		if info.IsDir() {
			if dirCanBeSkipped(path) {
				return filepath.SkipDir
			}

			language, error := matchedProgrammingLanguage(path)
			if error == nil {
				dir, name := dirAndName(rootLocation, path)
				projectData := ProjectData{BaseDir: dir, FullName: name, Type: language}
				*projects = append(*projects, projectData)
			}
		}

		return nil
	}
}

var skipIndicators = []string{".git", ".terraform", "node_modules", ".venv", "venv", "target", "build"}

func dirCanBeSkipped(path string) bool {
	for _, dirName := range skipIndicators {
		if filepath.Base(path) == dirName {
			return true
		}
	}
	return false
}

var fileToLanguageMapping = map[string]string{
	"build.gradle":     "java",
	"build.gradle.kts": "kotlin",
	"pom.xml":          "java",
	"go.mod":           "go",
	"Pipfile":          "python",
	"requirements.txt": "python",
	"package.json":     "node",
}

func matchedProgrammingLanguage(path string) (string, error) {

	for file, language := range fileToLanguageMapping {
		indicatorFile := filepath.Join(path, file)
		_, err := os.Stat(indicatorFile)
		if err == nil {
			return language, nil
		}
	}

	indicatorFile := filepath.Join(path, "Makefile")
	_, err := os.Stat(indicatorFile)
	if err == nil {
		return "Any", nil
	}

	return "", errors.New("No language matched")
}

func dirAndName(rootLocation string, path string) (string, string) {
	if path == rootLocation {
		return filepath.Dir(path), filepath.Base(path)
	} else {
		return rootLocation, strings.Replace(path, rootLocation+"/", "", -1)
	}
}
