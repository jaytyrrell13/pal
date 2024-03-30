package cmd

import (
	"os"

	"github.com/jaytyrrell13/pal/cmd/add"
	"github.com/jaytyrrell13/pal/cmd/install"
	"github.com/jaytyrrell13/pal/cmd/list"
	"github.com/jaytyrrell13/pal/cmd/make"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: "0.3.0",
	Use:     "pal",
	Short:   "Helps manage the aliases for your projects",
}

func Execute() {
	rootCmd.AddCommand(add.AddCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(install.InstallCmd)
	rootCmd.AddCommand(make.MakeCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
