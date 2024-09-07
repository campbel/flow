/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/campbel/flow/flow"
	"github.com/campbel/flow/meta/config"
	"github.com/campbel/flow/types"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flow",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(flowfile types.Flowfile) {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Use: "no-help", Hidden: true})
	rootCmd.AddGroup(&cobra.Group{ID: "custom", Title: "Flows:"})
	for name, flow := range flowfile.Flows {
		rootCmd.AddCommand(createFlowCommand(flowfile, name, flow))
	}
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

func createFlowCommand(flowfile types.Flowfile, name string, fl types.Flow) *cobra.Command {
	return &cobra.Command{
		Use:     name,
		Short:   "A brief description of your command",
		Long:    `A longer description that spans multiple lines and likely contains`,
		GroupID: "custom",
		RunE: func(cmd *cobra.Command, args []string) error {
			return flow.NewContext(flowfile, config.Workdir, config.Homedir, strings.Join(args, " ")).Execute(cmd.Context(), fl)
		},
	}
}
