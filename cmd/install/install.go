package install

import (
	"os"

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
	Use:     "install",
	Short:   "Create the configuration file used by pal",
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		path := pathFlag
		editorCmd := editorCmdFlag

		if path == "" {
			path = prompts.StringPrompt("What is the path to your projects?", os.Stdin)
		}

		if editorCmd == "" {
			editorCmd = prompts.StringPrompt("What is the editor command?", os.Stdin)
		}

		AppFs := afero.NewOsFs()

		configDirPath, configDirPathErr := pkg.ConfigDirPath()
		cobra.CheckErr(configDirPathErr)

		if pkg.FileMissing(AppFs, configDirPath) {
			configDirErr := pkg.MakeConfigDir(AppFs)
			cobra.CheckErr(configDirErr)
		}

		configFilePath, configFilePathErr := pkg.ConfigFilePath()
		cobra.CheckErr(configFilePathErr)

		if pkg.FileMissing(AppFs, configFilePath) {
			c := pkg.Config{
				Path:      path,
				EditorCmd: editorCmd,
			}

			json, jsonErr := c.AsJson()
			cobra.CheckErr(jsonErr)

			writeFileErr := pkg.WriteFile(AppFs, configFilePath, json, 0o644)
			cobra.CheckErr(writeFileErr)

			return
		}

		configFile, readConfigFileErr := pkg.ReadFile(AppFs, configFilePath)
		cobra.CheckErr(readConfigFileErr)

		c, configFileErr := pkg.FromJson(configFile)
		cobra.CheckErr(configFileErr)

		if c.Path != path {
			c.Path = path
		}

		if c.EditorCmd != editorCmd {
			c.EditorCmd = editorCmd
		}

		json, jsonErr := c.AsJson()
		cobra.CheckErr(jsonErr)

		writeFileErr := pkg.WriteFile(AppFs, configFilePath, json, 0o644)
		cobra.CheckErr(writeFileErr)
	},
}

func init() {
	InstallCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to to your projects")
	InstallCmd.Flags().StringVarP(&editorCmdFlag, "editorCmd", "e", "", "Editor command e.g. (nvim, subl, code)")
}
