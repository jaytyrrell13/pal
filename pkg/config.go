package pkg

import (
	"encoding/json"
	"os"

	"github.com/spf13/afero"
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

func MakeConfigDir(fs afero.Fs) error {
	return fs.MkdirAll(ConfigDirPath(), 0o750)
}

func FromJson(j []byte) (Config, error) {
	var c Config
	unmarshalErr := json.Unmarshal(j, &c)

	return c, unmarshalErr
}

func SaveExtraDir(path string) error {
	AppFs := afero.NewOsFs()
	configFile, readConfigFileErr := ReadFile(AppFs, ConfigFilePath())
	if readConfigFileErr != nil {
		return readConfigFileErr
	}

	c, configFileErr := FromJson(configFile)
	if configFileErr != nil {
		return configFileErr
	}

	c.Extras = append(c.Extras, path)
	saveErr := c.Save()
	if saveErr != nil {
		return saveErr
	}

	return nil
}

func (c Config) Save() error {
	json, jsonErr := c.AsJson()
	if jsonErr != nil {
		return jsonErr
	}

	writeFileErr := os.WriteFile(ConfigFilePath(), json, 0o644)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}

func (c Config) AsJson() ([]byte, error) {
	return json.Marshal(c)
}
