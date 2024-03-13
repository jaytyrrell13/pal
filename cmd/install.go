package cmd

import (
	"os"

	"github.com/jaytyrrell13/pal/pkg/config"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Create the configuration file used by pal",
	Run: func(cmd *cobra.Command, args []string) {
		path := prompts.StringPrompt("What is the path to your projects?")

		editorCmd := prompts.StringPrompt("What is the editor command?")

		if config.ConfigFileMissing() {
			c := config.Config{
				Path:      path,
				EditorCmd: editorCmd,
			}

			c.Save()

			return
		}

		configFile, openErr := os.ReadFile(config.ConfigFilePath())
		cobra.CheckErr(openErr)

		var c config.Config
		c = c.FromJson(configFile)

		if c.Path != path {
			c.Path = path
		}

		if c.EditorCmd != editorCmd {
			c.EditorCmd = editorCmd
		}

		c.Save()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
