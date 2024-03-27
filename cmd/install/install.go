package install

import (
	"os"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/cobra"
)

var (
	pathFlag      string
	editorCmdFlag string
)

var InstallCmd = &cobra.Command{
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

		if pkg.ConfigDirMissing() {
			pkg.MakeConfigDir()
		}

		if pkg.ConfigFileMissing() {
			c := pkg.Config{
				Path:      path,
				EditorCmd: editorCmd,
			}

			c.Save()

			return
		}

		configFile, openErr := os.ReadFile(pkg.ConfigFilePath())
		cobra.CheckErr(openErr)

		c := pkg.FromJson(configFile)

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
	InstallCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to to your projects")
	InstallCmd.Flags().StringVarP(&editorCmdFlag, "editorCmd", "e", "", "Editor command e.g. (nvim, subl, code)")
}
