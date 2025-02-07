package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/jaytyrrell13/pal/internal/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type InstallPrompts struct {
	shell string
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()
		ip, err := RunPrompts()
		if err != nil {
			return err
		}

		return RunInstallCmd(fs, ip)
	},
}

func RunPrompts() (InstallPrompts, error) {
	shell, err := ui.Select("What shell do you use?", []huh.Option[string]{
		huh.NewOption("Bash", "bash"),
		huh.NewOption("Fish", "fish"),
		huh.NewOption("Zsh", "zsh"),
	})
	if err != nil {
		return InstallPrompts{}, err
	}

	return InstallPrompts{
		shell: shell,
	}, nil
}

func RunInstallCmd(fs afero.Fs, ip InstallPrompts) error {
	configDirPath, configDirPathErr := config.ConfigDirPath()
	if configDirPathErr != nil {
		return configDirPathErr
	}

	mkDirErr := fs.MkdirAll(configDirPath, 0o750)
	if mkDirErr != nil {
		return mkDirErr
	}

	configFilePath, configFilePathErr := config.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	_, statErr := fs.Stat(configFilePath)
	if !errors.Is(statErr, os.ErrNotExist) {
		return errors.New("Config file already exists.")
	}

	c := config.NewConfig(ip.shell)

	j, jsonErr := json.Marshal(c)
	if jsonErr != nil {
		return jsonErr
	}

	writeFileErr := afero.WriteFile(fs, configFilePath, j, 0o644)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
}
