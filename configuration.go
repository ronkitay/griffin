package main

import (
	"encoding/json"
	"os"
)

type UserConfiguration struct {
	RepoRoots []string `json:"repoRoots"`
}

type Configuration struct {
	repoListLocation  string
	userConfiguration UserConfiguration
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

func loadConfiguration() Configuration {
	configurationDirectory := os.Getenv("HOME") + "/.config/rr"
	repoListLocation := configurationDirectory + "/rr.list"

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

	return Configuration{repoListLocation: repoListLocation, userConfiguration: userConfiguration}
}
