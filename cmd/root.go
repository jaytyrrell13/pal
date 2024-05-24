package cmd

import (
	"os"

	"github.com/jaytyrrell13/pal/cmd/add"
	"github.com/jaytyrrell13/pal/cmd/clean"
	"github.com/jaytyrrell13/pal/cmd/install"
	"github.com/jaytyrrell13/pal/cmd/list"
	"github.com/jaytyrrell13/pal/cmd/make"
	"github.com/jaytyrrell13/pal/cmd/refresh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pal",
	Short: "Helps manage the aliases for your projects",
}

func Execute(version string) {
	rootCmd.Version = version
	rootCmd.AddCommand(add.AddCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(install.InstallCmd)
	rootCmd.AddCommand(make.MakeCmd)
	rootCmd.AddCommand(clean.CleanCmd)
	rootCmd.AddCommand(refresh.RefreshCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
