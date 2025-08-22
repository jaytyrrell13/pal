package config

import (
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage pal config",
	}

	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewUpdateCmd())

	return cmd
}
