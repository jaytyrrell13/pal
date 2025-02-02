package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/internal"
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
	var shell string
	err := huh.NewSelect[string]().
		Title("What shell do you use?").
		Options(
			huh.NewOption("Bash", "bash"),
			huh.NewOption("Fish", "fish"),
			huh.NewOption("Zsh", "zsh"),
		).
		Value(&shell).
		WithTheme(huh.ThemeBase()).
		Run()
	if err != nil {
		return InstallPrompts{}, err
	}

	return InstallPrompts{
		shell: shell,
	}, nil
}

func RunInstallCmd(fs afero.Fs, ip InstallPrompts) error {
	configDirPath, configDirPathErr := internal.ConfigDirPath()
	if configDirPathErr != nil {
		return configDirPathErr
	}

	mkDirErr := fs.MkdirAll(configDirPath, 0o750)
	if mkDirErr != nil {
		return mkDirErr
	}

	configFilePath, configFilePathErr := internal.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	_, statErr := fs.Stat(configFilePath)
	if !errors.Is(statErr, os.ErrNotExist) {
		return errors.New("Config file already exists.")
	}

	c := internal.NewConfig(ip.shell)

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
