package config

import (
	"errors"
	"io"
	"os"
	"reflect"

	cfg "github.com/jaytyrrell13/pal/internal/config"
	"github.com/jaytyrrell13/pal/internal/messages"
	"github.com/jaytyrrell13/pal/internal/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List pal config",
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
}

func CheckListPrerequisites(fs afero.Fs) error {
	configFileExists, configFileExistsErr := cfg.ConfigFileExists(fs)
	if configFileExistsErr != nil {
		return configFileExistsErr
	}

	if !configFileExists {
		return errors.New(messages.Errors["configMissing"])
	}

	return nil
}

func RunListCmd(fs afero.Fs, w io.Writer) error {
	c, configErr := cfg.ReadConfigFile(fs)
	if configErr != nil {
		return configErr
	}

	headers := []string{"Key", "Value"}
	rows := [][]string{}

	v := reflect.ValueOf(c)
	typeOfS := v.Type()

	for i := range v.NumField() {
		fieldName := typeOfS.Field(i).Name

		if fieldName != "Aliases" {
			rows = append(rows, []string{fieldName, v.Field(i).String()})
		}
	}

	t := ui.Table(headers, rows)

	_, writeErr := w.Write([]byte(t.String()))

	return writeErr
}
