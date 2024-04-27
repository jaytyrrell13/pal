package make

import (
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

		configFilePath, configFilePathErr := pkg.ConfigFilePath()
		cobra.CheckErr(configFilePathErr)

		if pkg.FileMissing(AppFs, configFilePath) {
			cobra.CheckErr("Config file does not exist. Please run install command first.")
		}

		c, readConfigFileErr := pkg.ReadFile(AppFs, configFilePath)
		cobra.CheckErr(readConfigFileErr)

		jsonConfig, fromJsonErr := pkg.FromJson(c)
		cobra.CheckErr(fromJsonErr)

		path := jsonConfig.Path

		if strings.HasPrefix(path, "~/") {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, path[2:])
		}

		files, readDirErr := afero.ReadDir(AppFs, path)
		cobra.CheckErr(readDirErr)

		var projectPaths []string
		for _, file := range files {
			if file.Name() != ".DS_Store" {
				projectPaths = append(projectPaths, jsonConfig.Path+"/"+file.Name())
			}
		}

		projectPaths = append(projectPaths, jsonConfig.Extras...)

		var output string
		for _, path := range projectPaths {
			alias := prompts.StringPrompt(fmt.Sprintf("Alias for (%s) Leave blank to skip.", path), os.Stdin)

			if alias == "" {
				continue
			}

			output += pkg.MakeAliasCommands(alias, path, jsonConfig)
		}

		if output == "" {
			return
		}

		aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
		cobra.CheckErr(aliasFilePathErr)

		writeErr := pkg.WriteFile(AppFs, aliasFilePath, []byte(output), 0o755)
		cobra.CheckErr(writeErr)

		fmt.Println("\nDon't forget to source ~/.pal file in your shell!")
	},
}
