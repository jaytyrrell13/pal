package config

import (
	"github.com/jaytyrrell13/pal/cmd/config/list"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
}

func init() {
	ConfigCmd.AddCommand(list.ConfigListCmd)
}
