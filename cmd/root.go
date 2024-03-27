package cmd

import (
	"os"

	"github.com/jaytyrrell13/pal/cmd/add"
	"github.com/jaytyrrell13/pal/cmd/aliases"
	"github.com/jaytyrrell13/pal/cmd/install"
	"github.com/jaytyrrell13/pal/cmd/make"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pal",
	Short: "Helps manage the aliases for your projects",
}

func Execute() {
	rootCmd.AddCommand(add.AddCmd)
	rootCmd.AddCommand(aliases.AliasesCmd)
	rootCmd.AddCommand(install.InstallCmd)
	rootCmd.AddCommand(make.MakeCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
