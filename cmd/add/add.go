package add

import (
	"fmt"
	"os"

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
		fmt.Println("~/.pal file is missing. Running make command now.")

		makeCmdErr := make.RunMakeCmd()
		if makeCmdErr != nil {
			return makeCmdErr
		}
	}

	name := nameFlag
	path := pathFlag

	if name == "" {
		name = prompts.StringPrompt("What is the name of the alias?", os.Stdin)
	}

	if path == "" {
		path = prompts.StringPrompt("What is the path for the alias?", os.Stdin)
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
