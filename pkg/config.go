package pkg

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Path      string
	EditorCmd string
	Extras    []string
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

func ConfigDirMissing() bool {
	_, e := os.Stat(ConfigDirPath())
	return errors.Is(e, os.ErrNotExist)
}

func ConfigFileMissing() bool {
	_, e := os.Stat(ConfigFilePath())
	return errors.Is(e, os.ErrNotExist)
}

func ReadConfigFile() Config {
	configFile, openErr := os.ReadFile(ConfigFilePath())
	cobra.CheckErr(openErr)

	return FromJson(configFile)
}

func FromJson(j []byte) Config {
	var c Config
	unmarshalErr := json.Unmarshal(j, &c)
	cobra.CheckErr(unmarshalErr)

	return c
}

func SaveExtraDir(path string) {
	configFile, openErr := os.ReadFile(ConfigFilePath())
	cobra.CheckErr(openErr)

	c := FromJson(configFile)

	c.Extras = append(c.Extras, path)
	c.Save()
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
