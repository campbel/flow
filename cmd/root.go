/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/campbel/flow/flow"
	"github.com/campbel/flow/meta/config"
	"github.com/campbel/flow/meta/logger"
	"github.com/campbel/flow/types"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flow",
	Short: "",
	Long:  ``,
}

func Execute(flowfile types.Flowfile) {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Use: "no-help", Hidden: true})
	rootCmd.AddGroup(&cobra.Group{ID: "custom", Title: "Flows:"})
	for name, flow := range flowfile.Flows {
		rootCmd.AddCommand(createFlowCommand(flowfile, name, flow))
	}

	if err := rootCmd.Execute(); err != nil {
		if exitError, ok := err.(*types.ExitError); ok {
			logger.Error(exitError.Error(), "code", exitError.ExitCode())
			os.Exit(exitError.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}

func init() {
}

func createFlowCommand(flowfile types.Flowfile, name string, fl types.Flow) *cobra.Command {
	return &cobra.Command{
		Use:     name,
		Short:   "",
		Long:    ``,
		GroupID: "custom",
		RunE: func(cmd *cobra.Command, args []string) error {
			return flow.NewContext(flowfile, config.Workdir, config.Homedir, strings.Join(args, " ")).Execute(cmd.Context(), fl)
		},
	}
}
