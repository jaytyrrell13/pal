package cmd

import (
	"github.com/jaytyrrell13/pal/cmd/add"
	"github.com/jaytyrrell13/pal/cmd/clean"
	"github.com/jaytyrrell13/pal/cmd/config"
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

func Execute(version string) error {
	rootCmd.Version = version
	rootCmd.AddCommand(add.AddCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(install.InstallCmd)
	rootCmd.AddCommand(make.MakeCmd)
	rootCmd.AddCommand(clean.CleanCmd)
	rootCmd.AddCommand(refresh.RefreshCmd)
	rootCmd.AddCommand(config.ConfigCmd)

	return rootCmd.Execute()
}
