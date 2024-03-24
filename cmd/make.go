package cmd

import (
	"fmt"
	"os"

	"github.com/jaytyrrell13/pal/pkg/config"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/cobra"
)

func getProjectPaths(config config.Config) []string {
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

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Create the aliases file",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.ReadConfigFile()

		if config.ConfigFileMissing() {
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

			output += fmt.Sprintf("alias %s=\"cd %s\"\n", alias, path)

			if c.EditorCmd != "" {
				output += fmt.Sprintf("alias %s=\"cd %s && %s\"\n", "e"+alias, path, c.EditorCmd)
			}
		}

		if output == "" {
			return
		}

		homedir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		writeFileErr := os.WriteFile(homedir+"/.pal", []byte(output), 0o755)
		cobra.CheckErr(writeFileErr)

		fmt.Println("\nDon't forget to source ~/.pal file in your shell!")
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)
}
