package make

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/cobra"
)

func getProjectPaths(config pkg.Config) []string {
	path := config.Path

	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}

	files, err := os.ReadDir(path)
	cobra.CheckErr(err)

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
		if pkg.ConfigFileMissing() {
			cobra.CheckErr("Config file does not exist. Please run install command first.")
		}

		c := pkg.ReadConfigFile()

		paths := getProjectPaths(c)
		var output string
		for _, path := range paths {
			alias := prompts.StringPrompt(fmt.Sprintf("Alias for (%s) Leave blank to skip.", path))

			if alias == "" {
				continue
			}

			output += pkg.MakeAliasCommands(alias, path, c)
		}

		if output == "" {
			return
		}

		writeFileErr := os.WriteFile(pkg.AliasFilePath(), []byte(output), 0o755)
		cobra.CheckErr(writeFileErr)

		fmt.Println("\nDon't forget to source ~/.pal file in your shell!")
	},
}
