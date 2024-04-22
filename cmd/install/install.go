package install

import (
	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/afero"
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

		AppFs := afero.NewOsFs()

		if pkg.FileMissing(AppFs, pkg.ConfigDirPath()) {
			configDirErr := pkg.MakeConfigDir(AppFs)
			cobra.CheckErr(configDirErr)
		}

		if pkg.FileMissing(AppFs, pkg.ConfigFilePath()) {
			c := pkg.Config{
				Path:      path,
				EditorCmd: editorCmd,
			}

			saveErr := c.Save()
			cobra.CheckErr(saveErr)

			return
		}

		configFile, readConfigFileErr := pkg.ReadFile(AppFs, pkg.ConfigFilePath())
		cobra.CheckErr(readConfigFileErr)

		c, configFileErr := pkg.FromJson(configFile)
		cobra.CheckErr(configFileErr)

		if c.Path != path {
			c.Path = path
		}

		if c.EditorCmd != editorCmd {
			c.EditorCmd = editorCmd
		}

		saveErr := c.Save()
		cobra.CheckErr(saveErr)
	},
}

func init() {
	InstallCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to to your projects")
	InstallCmd.Flags().StringVarP(&editorCmdFlag, "editorCmd", "e", "", "Editor command e.g. (nvim, subl, code)")
}
