package set

import (
	"errors"
	"reflect"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var ConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set pal configs",
	RunE: func(_ *cobra.Command, args []string) error {
		return RunConfigSetCmd(args)
	},
}

func RunConfigSetCmd(args []string) error {
	AppFs := afero.NewOsFs()

	if len(args) != 2 {
		return errors.New("requires 2 args")
	}

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	jsonConfig, readConfigFileErr := pkg.ReadFile(AppFs, configFilePath)
	if readConfigFileErr != nil {
		return readConfigFileErr
	}

	c, fromJsonErr := pkg.FromJson(jsonConfig)
	if fromJsonErr != nil {
		return fromJsonErr
	}

	values := reflect.ValueOf(c)
	types := values.Type()

	key := args[0]
	value := args[1]

	titleCasedKey := cases.Title(language.English).String(key)

	_, ok := types.FieldByName(titleCasedKey)
	if !ok {
		return errors.New("Key must exist in Config struct")
	}

	switch titleCasedKey {
	case "Path":
		c.Path = value
	case "EditorCmd":
		c.EditorCmd = value
	case "Shell":
		c.Shell = value
	case "Extras":
		c.Extras = append(c.Extras, value)
	}

	json, jsonErr := c.AsJson()
	if jsonErr != nil {
		return jsonErr
	}

	writeFileErr := pkg.WriteFile(AppFs, configFilePath, json, 0o644)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}
