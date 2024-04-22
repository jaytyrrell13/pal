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

func getProjectPaths(config pkg.Config) []string {
	path := config.Path

	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}

	files, readDirErr := os.ReadDir(path)
	cobra.CheckErr(readDirErr)

	var projectPaths []string
	for _, file := range files {
		if file.Name() != ".DS_Store" {
			projectPaths = append(projectPaths, config.Path+"/"+file.Name())
		}
	}

	projectPaths = append(projectPaths, config.Extras...)

	return projectPaths
}

var MakeCmd = &cobra.Command{
	Use:   "make",
	Short: "Create the aliases file",
	Run: func(cmd *cobra.Command, args []string) {
		RunMakeCmd()
	},
}

func RunMakeCmd() {
	AppFs := afero.NewOsFs()
	if pkg.FileMissing(AppFs, pkg.ConfigFilePath()) {
		cobra.CheckErr("Config file does not exist. Please run install command first.")
	}

	c, readConfigFileErr := pkg.ReadFile(AppFs, pkg.ConfigFilePath())
	cobra.CheckErr(readConfigFileErr)

	jsonConfig, fromJsonErr := pkg.FromJson(c)
	cobra.CheckErr(fromJsonErr)

	paths := getProjectPaths(jsonConfig)
	var output string
	for _, path := range paths {
		alias := prompts.StringPrompt(fmt.Sprintf("Alias for (%s) Leave blank to skip.", path))

		if alias == "" {
			continue
		}

		output += pkg.MakeAliasCommands(alias, path, jsonConfig)
	}

	if output == "" {
		return
	}

	writeErr := pkg.WriteFile(AppFs, pkg.AliasFilePath(), []byte(output), 0o755)
	cobra.CheckErr(writeErr)

	fmt.Println("\nDon't forget to source ~/.pal file in your shell!")
}
