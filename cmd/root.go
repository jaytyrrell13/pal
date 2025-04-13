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
	rootCmd.AddCommand(NewCreateCmd())
	rootCmd.AddCommand(NewInstallCmd())
	rootCmd.AddCommand(NewListCmd())
	rootCmd.AddCommand(NewRemoveCmd())
	rootCmd.AddCommand(NewUpdateCmd())
	rootCmd.AddCommand(config.NewConfigCmd())

	return rootCmd.Execute()
}
