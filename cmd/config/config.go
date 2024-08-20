package config

import (
	"github.com/jaytyrrell13/pal/cmd/config/list"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage pal config",
}

func init() {
	ConfigCmd.AddCommand(list.ConfigListCmd)
}
