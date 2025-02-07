package config

import (
	"os"

	"github.com/jaytyrrell13/pal/internal/alias"
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
