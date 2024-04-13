package configuration

import (
	"encoding/json"
	"os"
)

type IdeConfiguration struct {
	DefaultIDE string `json:"default"`
	GoLang     string `json:"go"`
	Java       string `json:"java"`
	Kotlin     string `json:"kotlin"`
	Python     string `json:"python"`
	NodeJS     string `json:"node"`
}

type UserConfiguration struct {
	RepoRoots        []string         `json:"repoRoots"`
	IdeConfiguration IdeConfiguration `json:"ideConfiguration"`
}

type Configuration struct {
	RepoListLocation    string
	ProjectListLocation string
	UserConfiguration   UserConfiguration
}

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// File does not exist
		return false, nil
	} else if err != nil {
		// An error occurred (other than file not existing)
		return false, err
	}

	// File exists
	return true, nil
}

func LoadConfiguration() Configuration {
	configurationDirectory := os.Getenv("HOME") + "/.config/griffin"
	repoListLocation := configurationDirectory + "/repo.list"
	projectListLocation := configurationDirectory + "/project.list"

	var userConfiguration UserConfiguration

	if exists, _ := fileExists(configurationDirectory + "/config.json"); exists == true {
		jsonConfigFile, fileOpenError := os.ReadFile(configurationDirectory + "/config.json")
		if fileOpenError != nil {
			os.Exit(1)
		}

		jsonReadError := json.Unmarshal(jsonConfigFile, &userConfiguration)
		if jsonReadError != nil {
			panic("Error reading JSON Configuration:" + jsonReadError.Error())
		}
	} else {
		userConfiguration = UserConfiguration{}
	}

	return Configuration{RepoListLocation: repoListLocation, ProjectListLocation: projectListLocation, UserConfiguration: userConfiguration}
}
