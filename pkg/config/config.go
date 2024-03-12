package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Path      string
	EditorCmd string
}

func ConfigDirPath() string {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return homeDir + "/.config/pal"
}

func ConfigFilePath() string {
	return ConfigDirPath() + "/config.json"
}

func MakeConfigDir() {
	pathErr := os.MkdirAll(ConfigDirPath(), 0o750)
	cobra.CheckErr(pathErr)
}

func ConfigFileExists() bool {
	_, e := os.Stat(ConfigFilePath())
	return !errors.Is(e, os.ErrNotExist)
}

func WriteConfigFile(b []byte) {
	writeFileErr := os.WriteFile(ConfigFilePath(), b, 0o644)
	cobra.CheckErr(writeFileErr)
}

func (c Config) FromJson(j []byte) Config {
	unmarshalErr := json.Unmarshal(j, &c)
	cobra.CheckErr(unmarshalErr)

	return c
}

func (c Config) Save() {
	writeFileErr := os.WriteFile(ConfigFilePath(), c.AsJson(), 0o644)
	cobra.CheckErr(writeFileErr)
}

func (c Config) AsJson() []byte {
	b, marshalErr := json.Marshal(c)
	cobra.CheckErr(marshalErr)

	return b
}
