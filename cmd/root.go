package cmd

import (
	"github.com/jaytyrrell13/pal/cmd/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "pal",
	Short:         "Command line tool to create, update, and remove shell aliases.",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute(version string) error {
	rootCmd.Version = version
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(config.ConfigCmd)

	return rootCmd.Execute()
}
