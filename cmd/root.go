package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          binaryName,
	Short:        "An example CLI demonstrating Cobra, Viper, XDG, and colored output",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if shouldSkipConfigLoad(cmd) {
			return nil
		}

		return loadConfig(cmd)
	},
}

// Execute runs the root command while keeping the main package free of CLI
// wiring details. When the binary moves beyond the current demo objective,
// replace or add command registrations inside package cmd and leave main as a
// thin entrypoint.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		fmt.Sprintf("path to a TOML config file (default XDG path: %s)", defaultUserConfigFilePath()),
	)
}
