package install

import (
	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	pathFlag      string
	editorCmdFlag string
	shell         string
)

var InstallCmd = &cobra.Command{
	Use:     "install",
	Short:   "Create the configuration file used by pal",
	Aliases: []string{"i"},
	RunE: func(_ *cobra.Command, _ []string) error {
		return RunInstallCmd()
	},
}

func init() {
	InstallCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to your projects")
	InstallCmd.Flags().StringVarP(&editorCmdFlag, "editorCmd", "e", "", "Editor command e.g. (nvim, subl, code)")
	InstallCmd.Flags().StringVarP(&shell, "shell", "s", "", "Your interactive shell e.g. (bash/zsh, fish)")
}

func RunInstallCmd() error {
	path := pathFlag
	editorCmd := editorCmdFlag

	if path == "" {
		pathString, pathErr := prompts.Input("What is the path to your projects?", "/Users/john/Code")

		if pathErr != nil {
			return pathErr
		}

		path = pathString
	}

	if editorCmd == "" {
		editorCmdString, editorCmdErr := prompts.Input("What is the editor command?", "nvim, subl, code")

		if editorCmdErr != nil {
			return editorCmdErr
		}

		editorCmd = editorCmdString
	}

	if shell == "" {
		shellString, shellErr := prompts.Select("What shell do you use?", huh.NewOptions("Bash", "ZSH", "Fish"))

		if shellErr != nil {
			return shellErr
		}

		shell = shellString
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
			Shell:     shell,
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

	if c.Shell != shell {
		c.Shell = shell
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
