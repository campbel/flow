/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"path/filepath"

	"github.com/campbel/flow/cmd"
	"github.com/campbel/flow/meta/config"
	"github.com/campbel/flow/types"
)

func main() {
	flowFile, err := types.LoadFlowfile(filepath.Join(config.Workdir, "flow.yaml"))
	if err != nil {
		panic(err)
	}

	cmd.Execute(flowFile)
}
