package install

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	pathFlag       string
	editorCmdFlag  string
	editorModeFlag string
	shellFlag      string
)

var InstallCmd = &cobra.Command{
	Use:     "install",
	Short:   "Create the configuration file used by pal",
	Aliases: []string{"i"},
	RunE: func(_ *cobra.Command, _ []string) error {
		appFs := afero.NewOsFs()

		return RunInstallCmd(appFs)
	},
}

func init() {
	InstallCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to your projects")
	InstallCmd.Flags().StringVarP(&editorModeFlag, "editorMode", "", "", "Editor command mode e.g. (skip, same, unique)")
	InstallCmd.Flags().StringVarP(&editorCmdFlag, "editorCmd", "e", "", "Editor command e.g. (nvim, subl, code)")
	InstallCmd.Flags().StringVarP(&shellFlag, "shell", "s", "", "Your interactive shell e.g. (bash, zsh, fish)")
}

func RunInstallCmd(appFs afero.Fs) error {
	path := pathFlag
	editorMode := editorModeFlag
	editorCmd := editorCmdFlag
	shell := shellFlag

	if path == "" {
		pathString, pathErr := ui.Input("What is the path to your projects?", "/Users/john/Code")

		if pathErr != nil {
			return pathErr
		}

		path = pathString
	}

	if pkg.FileMissing(appFs, path) {
		return fmt.Errorf("Path '%s' does not exist. Please try again.", path)
	}

	if editorMode == "" {
		options := []huh.Option[string]{
			huh.NewOption("Skip", "skip"),
			huh.NewOption("Same", "same"),
			huh.NewOption("Unique", "unique"),
		}
		editorModeString, editorModeErr := ui.Select("How would you like to use the editor command?", options)

		if editorModeErr != nil {
			return editorModeErr
		}

		editorMode = editorModeString
	}

	if editorCmd == "" {
		editorCmdString, editorCmdErr := ui.Input("What is the editor command?", "nvim, subl, code")

		if editorCmdErr != nil {
			return editorCmdErr
		}

		editorCmd = editorCmdString
	}

	if shell == "" {
		options := []huh.Option[string]{
			huh.NewOption("Bash", pkg.BashShell),
			huh.NewOption("Zsh", pkg.ZshShell),
			huh.NewOption("Fish", pkg.FishShell),
		}
		shellString, shellErr := ui.Select("What shell do you use?", options)

		if shellErr != nil {
			return shellErr
		}

		shell = shellString
	}

	configDirPath, configDirPathErr := pkg.ConfigDirPath()
	if configDirPathErr != nil {
		return configDirPathErr
	}

	if pkg.FileMissing(appFs, configDirPath) {
		configDirErr := pkg.MakeConfigDir(appFs)
		if configDirErr != nil {
			return configDirErr
		}
	}

	if ok, _ := pkg.ConfigFileMissing(appFs); !ok {
		return nil
	}

	c, configErr := pkg.NewConfig(path, editorMode, editorCmd, shell)
	if configErr != nil {
		return configErr
	}

	return c.Save(appFs)
}
