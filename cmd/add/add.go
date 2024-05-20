package add

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/cmd/make"
	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/prompts"
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
	Run: func(cmd *cobra.Command, args []string) {
		err := RunAddCmd()
		cobra.CheckErr(err)
	},
}

func init() {
	AddCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "Name of the additional alias")
	AddCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to your additional directory")
}

func RunAddCmd() error {
	AppFs := afero.NewOsFs()

	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	if pkg.FileMissing(AppFs, aliasFilePath) {
		var runMake bool
		confirmErr := huh.NewConfirm().
			Title("Alias file is missing. Would you like to run make command now?").
			Value(&runMake).
			Affirmative("Yes").
			Negative("No").
			Run()

		if confirmErr != nil {
			return confirmErr
		}

		if runMake {
			fmt.Println("Running make command.")

			makeCmdErr := make.RunMakeCmd()
			if makeCmdErr != nil {
				return makeCmdErr
			}
		} else {
			fmt.Println("Please run `pal make` command manually.")
			return nil
		}
	}

	name := nameFlag
	path := pathFlag

	if name == "" {
		nameString, nameErr := prompts.Input("What is the name of the alias?", "foo")

		if nameErr != nil {
			return nameErr
		}

		name = nameString
	}

	if path == "" {
		pathString, pathErr := prompts.Input("What is the path for the alias?", "/Users/john/Documents")

		if pathErr != nil {
			return pathErr
		}

		path = pathString
	}

	saveExtraDirErr := pkg.SaveExtraDir(AppFs, path)
	if saveExtraDirErr != nil {
		return saveExtraDirErr
	}

	aliasesFile, openAliasesFileErr := os.OpenFile(aliasFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
	if openAliasesFileErr != nil {
		return openAliasesFileErr
	}

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	c, readConfigFileErr := pkg.ReadFile(AppFs, configFilePath)
	if readConfigFileErr != nil {
		return readConfigFileErr
	}

	jsonConfig, fromJsonErr := pkg.FromJson(c)
	if fromJsonErr != nil {
		return fromJsonErr
	}

	output := pkg.MakeAliasCommands(name, path, jsonConfig)

	if _, err := aliasesFile.Write([]byte(output)); err != nil {
		aliasesFile.Close()
		return err
	}

	if err := aliasesFile.Close(); err != nil {
		return err
	}

	return nil
}
