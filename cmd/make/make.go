package make

import (
	"fmt"
	"os"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/cobra"
)

func getProjectPaths(config pkg.Config) []string {
	files, err := os.ReadDir(config.Path)
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
		c := pkg.ReadConfigFile()

		if pkg.ConfigFileMissing() {
			fmt.Println("Config file does not exist. Please run install command first.")
			os.Exit(1)
		}

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
