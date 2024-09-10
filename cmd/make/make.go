package make

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

		if runInstall {
			fmt.Println("Running install command.")

			installCmdErr := install.RunInstallCmd(appFs)
			if installCmdErr != nil {
				return installCmdErr
			}
		} else {
			fmt.Println("Please run `pal install` command manually.")
			return nil
		}
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return homeErr
	}

	c, readConfigErr := pkg.ReadConfigFile(appFs)
	if readConfigErr != nil {
		return readConfigErr
	}

	path := c.Path

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	}

	files, readDirErr := afero.ReadDir(appFs, path)
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

	var aliases []pkg.Alias
	for _, path := range projectPaths {
		alias, aliasErr := ui.Input(fmt.Sprintf("Alias for (%s) Leave blank to skip.", path), "foo")

		if aliasErr != nil {
			return aliasErr
		}

		if alias == "" {
			continue
		}

		aliases = append(aliases, pkg.NewAlias(alias, path))
	}

	saveErr := pkg.SaveAliases(appFs, aliases, c)
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
	case pkg.Shell_Bash:
		return sourceBashFile(appFs, home)
	case pkg.Shell_Zsh:
		return sourceZshFile(appFs, home)
	case pkg.Shell_Fish:
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
