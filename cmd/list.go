package cmd

import (
	"errors"
	"io"
	"os"

	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/jaytyrrell13/pal/internal/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Display list of aliases",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		prereqsErr := CheckListPrerequisites(fs)
		if prereqsErr != nil {
			return prereqsErr
		}

		return RunListCmd(fs, os.Stdout)
	},
}

func CheckListPrerequisites(fs afero.Fs) error {
	configFileExists, configFileExistsErr := config.ConfigFileExists(fs)
	if configFileExistsErr != nil {
		return configFileExistsErr
	}

	if !configFileExists {
		return errors.New("Config file does not exist")
	}

	return nil
}

func RunListCmd(fs afero.Fs, w io.Writer) error {
	c, configErr := config.ReadConfigFile(fs)
	if configErr != nil {
		return configErr
	}

	headers := []string{"Alias", "Command"}
	rows := [][]string{}

	for _, a := range c.Aliases {
		rows = append(rows, []string{a.Name, a.Command})
	}

	t := ui.Table(headers, rows)

	_, writeErr := w.Write([]byte(t.String()))
	if writeErr != nil {
		return writeErr
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
