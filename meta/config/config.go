package config

import (
	"os"
)

var (
	DebugLogLevel = os.Getenv("DEBUG") == "true"
	Workdir       = loadWorkdir()
	Homedir       = loadHomedir()
)

func loadWorkdir() string {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return workdir
}

func loadHomedir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homedir
}
