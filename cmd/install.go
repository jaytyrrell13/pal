package cmd

import (
	"os"

	"github.com/jaytyrrell13/pal/pkg/config"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/cobra"
)

var (
	pathFlag      string
	editorCmdFlag string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Create the configuration file used by pal",
	Run: func(cmd *cobra.Command, args []string) {
		path := pathFlag
		editorCmd := editorCmdFlag

		if path == "" {
			path = prompts.StringPrompt("What is the path to your projects?")
		}

		if editorCmd == "" {
			editorCmd = prompts.StringPrompt("What is the editor command?")
		}

		if config.ConfigDirMissing() {
			config.MakeConfigDir()
		}

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

		c := config.FromJson(configFile)

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

	installCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to to your projects")
	installCmd.Flags().StringVarP(&editorCmdFlag, "editorCmd", "e", "", "Editor command e.g. (nvim, subl, code)")
}
