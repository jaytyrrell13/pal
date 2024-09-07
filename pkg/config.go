package pkg

import (
	"encoding/json"
	"os"

	"github.com/spf13/afero"
)

type Config struct {
	Path      string
	EditorCmd string
	Shell     string
	Extras    []string
}

func NewConfig(path string, editorCmd string, shell string) Config {
	return Config{
		Path:      path,
		EditorCmd: editorCmd,
		Shell:     shell,
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

func MakeConfigDir(fs afero.Fs) error {
	path, err := ConfigDirPath()
	if err != nil {
		return err
	}

	return fs.MkdirAll(path, 0o750)
}

func FromJson(j []byte) (Config, error) {
	var c Config
	unmarshalErr := json.Unmarshal(j, &c)

	return c, unmarshalErr
}

func SaveExtraDir(appFs afero.Fs, path string) error {
	configFilePath, configFilePathErr := ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	configFile, readConfigFileErr := ReadFile(appFs, configFilePath)
	if readConfigFileErr != nil {
		return readConfigFileErr
	}

	c, configFileErr := FromJson(configFile)
	if configFileErr != nil {
		return configFileErr
	}

	c.Extras = append(c.Extras, path)

	return c.Save(appFs)
}

func (c Config) AsJson() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Config) Save(appFs afero.Fs) error {
	configFilePath, configFilePathErr := ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	json, jsonErr := c.AsJson()
	if jsonErr != nil {
		return jsonErr
	}

	writeFileErr := WriteFile(appFs, configFilePath, json, 0o644)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}
