package make

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var MakeCmd = &cobra.Command{
	Use:   "make",
	Short: "Create the aliases file",
	Run: func(cmd *cobra.Command, args []string) {
		AppFs := afero.NewOsFs()

		err := RunMakeCmd(AppFs)
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func RunMakeCmd(fs afero.Fs) error {
	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	if pkg.FileMissing(fs, configFilePath) {
		return errors.New("Config file does not exist. Please run install command first.")
	}

	c, readConfigFileErr := pkg.ReadFile(fs, configFilePath)
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

	files, readDirErr := afero.ReadDir(fs, path)
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
		alias := prompts.StringPrompt(fmt.Sprintf("Alias for (%s) Leave blank to skip.", path))

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

	writeErr := pkg.WriteFile(fs, aliasFilePath, []byte(output), 0o755)
	if writeErr != nil {
		return writeErr
	}

	fmt.Println("\nDon't forget to source ~/.pal file in your shell!")

	return nil
}
