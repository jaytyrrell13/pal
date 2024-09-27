package set

import (
	"errors"
	"fmt"
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
		appFs := afero.NewOsFs()

		return RunConfigSetCmd(appFs, args)
	},
}

func RunConfigSetCmd(appFs afero.Fs, args []string) error {
	if len(args) != 2 {
		return errors.New("requires 2 args")
	}

	c, readConfigErr := pkg.ReadConfigFile(appFs)
	if readConfigErr != nil {
		return readConfigErr
	}

	values := reflect.ValueOf(c)
	types := values.Type()

	key := args[0]
	value := args[1]

	titleCasedKey := cases.Title(language.English).String(key)

	fmt.Println(titleCasedKey)

	_, ok := types.FieldByName(titleCasedKey)
	if !ok {
		return errors.New("Key must exist in Config struct")
	}

	switch titleCasedKey {
	case "Path":
		c.Path = value
	case "Editorcmd":
		c.Editorcmd = value
	case "Shell":
		if value != pkg.BashShell && value != pkg.ZshShell && value != pkg.FishShell {
			return fmt.Errorf("Shell must be either Bash, Zsh, or Fish. Received: %s", value)
		}

		c.Shell = value
	case "Extras":
		c.Extras = append(c.Extras, value)
	}

	return c.Save(appFs)
}
