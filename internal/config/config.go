package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/spf13/afero"
)

type Config struct {
	Shell   string
	Aliases []alias.Alias
}

func NewConfig(shell string) Config {
	return Config{
		Shell: shell,
	}
}

func ConfigDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()

	return homeDir + "/.config/pal", err
}

func ConfigFilePath() (string, error) {
	path, err := ConfigDirPath()
	if err != nil {
		return "", err
	}

	return path + "/config.json", nil
}

func ConfigFileExists(fs afero.Fs) (bool, error) {
	path, err := ConfigFilePath()
	if err != nil {
		return false, err
	}

	_, statErr := fs.Stat(path)
	return !errors.Is(statErr, os.ErrNotExist), nil
}

func ReadConfigFile(fs afero.Fs) (Config, error) {
	configFilePath, configFilePathErr := ConfigFilePath()
	if configFilePathErr != nil {
		return Config{}, configFilePathErr
	}

	b, readFileErr := afero.ReadFile(fs, configFilePath)
	if readFileErr != nil {
		return Config{}, readFileErr
	}

	var c Config
	jsonErr := json.Unmarshal(b, &c)
	if jsonErr != nil {
		return Config{}, jsonErr
	}

	return c, nil
}
