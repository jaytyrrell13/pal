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
		err := RunInstallCmd()
		cobra.CheckErr(err)
	},
}

func init() {
	InstallCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to to your projects")
	InstallCmd.Flags().StringVarP(&editorCmdFlag, "editorCmd", "e", "", "Editor command e.g. (nvim, subl, code)")
}

func RunInstallCmd() error {
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
	if configDirPathErr != nil {
		return configDirPathErr
	}

	if pkg.FileMissing(AppFs, configDirPath) {
		configDirErr := pkg.MakeConfigDir(AppFs)
		if configDirErr != nil {
			return configDirErr
		}
	}

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	if pkg.FileMissing(AppFs, configFilePath) {
		c := pkg.Config{
			Path:      path,
			EditorCmd: editorCmd,
		}

		json, jsonErr := c.AsJson()
		if jsonErr != nil {
			return jsonErr
		}

		writeFileErr := pkg.WriteFile(AppFs, configFilePath, json, 0o644)
		if writeFileErr != nil {
			return writeFileErr
		}

		return nil
	}

	configFile, readConfigFileErr := pkg.ReadFile(AppFs, configFilePath)
	if readConfigFileErr != nil {
		return readConfigFileErr
	}

	c, configFileErr := pkg.FromJson(configFile)
	if configFileErr != nil {
		return configFileErr
	}

	if c.Path != path {
		c.Path = path
	}

	if c.EditorCmd != editorCmd {
		c.EditorCmd = editorCmd
	}

	json, jsonErr := c.AsJson()
	if jsonErr != nil {
		return jsonErr
	}

	writeFileErr := pkg.WriteFile(AppFs, configFilePath, json, 0o644)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}
