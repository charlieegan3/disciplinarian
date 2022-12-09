/*
Copyright © 2022 Charlie Egan
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charlieegan3/disciplinarian/pkg/discipline"

	"github.com/charlieegan3/disciplinarian/pkg/config"

	"github.com/spf13/cobra"
)

var cfgFilePath string
var defaultCfgFilePath = ".disciplinarian.yaml"

var rootCmd = &cobra.Command{
	Use:          "disciplinarian",
	Short:        "A CLI tool to enforce policy rules on structured files",
	SilenceUsage: true, // don't show usage on err
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFilePath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		results, err := discipline.Run(cmd.Context(), cfg)
		if err != nil {
			return fmt.Errorf("failed to run discipline: %w", err)
		}

		if len(results) > 0 {
			fmt.Fprintf(os.Stdout, "%d violations\n\n", len(results))
			for _, r := range results {
				fmt.Fprintf(os.Stdout, "%s:\n", r.File)
				for _, m := range r.Messages {
					fmt.Fprintf(os.Stdout, "  - %s\n", m)
				}
			}
			os.Exit(1)
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", defaultCfgFilePath, "config file (default is $PWD/.disciplinarian.yaml)")
}

func initConfig() {}
