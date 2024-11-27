package make

import (
	"fmt"
	"os"

	"github.com/jaytyrrell13/pal/cmd/install"
	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var MakeCmd = &cobra.Command{
	Use:   "make",
	Short: "Create the aliases file",
	RunE: func(_ *cobra.Command, _ []string) error {
		appFs := afero.NewOsFs()

		return RunMakeCmd(appFs)
	},
}

func RunMakeCmd(appFs afero.Fs) error {
	if ok, _ := pkg.ConfigFileMissing(appFs); ok {
		runInstall, confirmErr := ui.Confirm("Config file does not exist. Would you like to run install command now?")

		if confirmErr != nil {
			return confirmErr
		}

		if !runInstall {
			fmt.Println("Please run `pal install` command manually.")
			return nil
		}

		fmt.Println("Running install command.")

		installCmdErr := install.RunInstallCmd(appFs)
		if installCmdErr != nil {
			return installCmdErr
		}
	}

	c, readConfigErr := pkg.ReadConfigFile(appFs)
	if readConfigErr != nil {
		return readConfigErr
	}

	files, readDirErr := afero.ReadDir(appFs, c.Path)
	if readDirErr != nil {
		return readDirErr
	}

	var projectPaths []string
	for _, file := range files {
		if file.Name() != ".DS_Store" {
			projectPaths = append(projectPaths, c.Path+"/"+file.Name())
		}
	}

	projectPaths = append(projectPaths, c.Extras...)

	var editorCmd string
	if c.Editormode == "same" {
		editorCmd = c.Editorcmd
	}

	var aliases []pkg.Alias
	for _, path := range projectPaths {
		alias, aliasErr := ui.Input(fmt.Sprintf("Alias for (%s) Leave blank to skip.", path), "foo")

		if aliasErr != nil {
			return aliasErr
		}

		if alias == "" {
			continue
		}

		if c.Editormode == "unique" {
			editorCmdString, editorCmdErr := ui.Input(fmt.Sprintf("What is the editor command for (%s)?", alias), "nvim, subl, code")

			if editorCmdErr != nil {
				return editorCmdErr
			}

			editorCmd = editorCmdString
		}

		aliases = append(aliases, pkg.NewAlias(alias, path, editorCmd))
	}

	saveErr := pkg.SaveAliases(appFs, aliases)
	if saveErr != nil {
		return saveErr
	}

	return sourceAliasFile(appFs, c)
}

func sourceAliasFile(appFs afero.Fs, config pkg.Config) error {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return homeErr
	}

	switch config.Shell {
	case pkg.BashShell:
		return sourceBashFile(appFs, home)
	case pkg.ZshShell:
		return sourceZshFile(appFs, home)
	case pkg.FishShell:
		return sourceFishFile(appFs, home)
	}

	return nil
}

func sourceBashFile(appFs afero.Fs, home string) error {
	data := []byte("\n[ -f \"$HOME/.config/pal/aliases\" ] && source \"$HOME/.config/pal/aliases\"")

	return pkg.AppendToFile(appFs, home+"/.bashrc", data)
}

func sourceZshFile(appFs afero.Fs, home string) error {
	data := []byte("\n[ -f \"$HOME/.config/pal/aliases\" ] && source \"$HOME/.config/pal/aliases\"")

	return pkg.AppendToFile(appFs, home+"/.zshrc", data)
}

func sourceFishFile(appFs afero.Fs, home string) error {
	data := []byte("if test -f " + home + "/.config/pal/aliases\n\tsource " + home + "/.config/pal/aliases\nend")

	return pkg.WriteFile(appFs, home+"/.config/fish/conf.d/pal.fish", data, 0o644)
}
