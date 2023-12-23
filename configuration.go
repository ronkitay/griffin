package main

import "os"

type Configuration struct {
	repoListLocation string
}

func loadConfiguration() Configuration {
	configurationDirectory := os.Getenv("HOME") + "/.config/rr"
	repoListLocation := configurationDirectory + "/rr.list"
	return Configuration{repoListLocation: repoListLocation}
}
