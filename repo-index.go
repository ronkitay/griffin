package main

import (
	"encoding/csv"
	"os"
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
		repoDir := parentDir + "/" + repoName

		switch locationType {
		case "dir":
			items = append(items, RepoData{repoDir: repoDir, repoName: repoName, locationType: "dir"})
		case "archive":
			fallthrough
		case "gitlab":
			fallthrough
		case "github":
			items = append(items, RepoData{repoDir: repoDir, repoName: repoName, url: url, locationType: locationType})
		}
	}
	return items
}

func loadIndexCsv() [][]string {
	configuration := loadConfiguration()

	file, fileOpenError := os.Open(configuration.repoListLocation)
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
