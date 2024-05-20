package make

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/cmd/install"
	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var MakeCmd = &cobra.Command{
	Use:   "make",
	Short: "Create the aliases file",
	Run: func(cmd *cobra.Command, args []string) {
		err := RunMakeCmd()
		cobra.CheckErr(err)
	},
}

func RunMakeCmd() error {
	AppFs := afero.NewOsFs()

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	if pkg.FileMissing(AppFs, configFilePath) {
		var runInstall bool
		confirmErr := huh.NewConfirm().
			Title("Config file does not exist. Would you like to run install command now?").
			Value(&runInstall).
			Affirmative("Yes").
			Negative("No").
			Run()

		if confirmErr != nil {
			return confirmErr
		}

		if runInstall {
			fmt.Println("Running install command.")

			installCmdErr := install.RunInstallCmd()
			if installCmdErr != nil {
				return installCmdErr
			}
		} else {
			fmt.Println("Please run `pal install` command manually.")
			return nil
		}
	}

	c, readConfigFileErr := pkg.ReadFile(AppFs, configFilePath)
	if readConfigFileErr != nil {
		return readConfigFileErr
	}

	jsonConfig, fromJsonErr := pkg.FromJson(c)
	if fromJsonErr != nil {
		return fromJsonErr
	}

	path := jsonConfig.Path

	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}

	files, readDirErr := afero.ReadDir(AppFs, path)
	if readDirErr != nil {
		return readDirErr
	}

	var projectPaths []string
	for _, file := range files {
		if file.Name() != ".DS_Store" {
			projectPaths = append(projectPaths, jsonConfig.Path+"/"+file.Name())
		}
	}

	projectPaths = append(projectPaths, jsonConfig.Extras...)

	var output string
	for _, path := range projectPaths {
		alias, aliasErr := prompts.Input(fmt.Sprintf("Alias for (%s) Leave blank to skip.", path), "foo")

		if aliasErr != nil {
			return aliasErr
		}

		if alias == "" {
			continue
		}

		output += pkg.MakeAliasCommands(alias, path, jsonConfig)
	}

	if output == "" {
		return nil
	}

	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	writeErr := pkg.WriteFile(AppFs, aliasFilePath, []byte(output), 0o755)
	if writeErr != nil {
		return writeErr
	}

	fmt.Println("Don't forget to source ~/.pal file in your shell!")

	return nil
}
