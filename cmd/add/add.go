package add

import (
	"fmt"

	"github.com/jaytyrrell13/pal/cmd/make"
	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	nameFlag string
	pathFlag string
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an alias for an additional directory",
	RunE: func(_ *cobra.Command, _ []string) error {
		appFs := afero.NewOsFs()

		return RunAddCmd(appFs, nameFlag, pathFlag)
	},
}

func init() {
	AddCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "Name of the additional alias")
	AddCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to your additional directory")
}

func RunAddCmd(appFs afero.Fs, name string, path string) error {
	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	if pkg.FileMissing(appFs, aliasFilePath) {
		runMake, confirmErr := ui.Confirm("Alias file is missing. Would you like to run make command now?")

		if confirmErr != nil {
			return confirmErr
		}

		if !runMake {
			fmt.Println("Please run `pal make` command manually.")
			return nil
		}

		fmt.Println("Running make command.")

		makeCmdErr := make.RunMakeCmd(appFs)
		if makeCmdErr != nil {
			return makeCmdErr
		}
	}

	if name == "" {
		nameString, nameErr := ui.Input("What is the name of the alias?", "foo")

		if nameErr != nil {
			return nameErr
		}

		name = nameString
	}

	if path == "" {
		pathString, pathErr := ui.Input("What is the path for the alias?", "/Users/john/Documents")

		if pathErr != nil {
			return pathErr
		}

		path = pathString
	}

	saveExtraDirErr := pkg.SaveExtraDir(appFs, path)
	if saveExtraDirErr != nil {
		return saveExtraDirErr
	}

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	c, readConfigFileErr := pkg.ReadFile(appFs, configFilePath)
	if readConfigFileErr != nil {
		return readConfigFileErr
	}

	jsonConfig, fromJsonErr := pkg.FromJson(c)
	if fromJsonErr != nil {
		return fromJsonErr
	}

	output := pkg.MakeAliasCommands(name, path, jsonConfig)

	return pkg.AppendToFile(appFs, aliasFilePath, []byte(output))
}
