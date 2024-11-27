package pkg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type Config struct {
	Path       string
	Editormode string
	Editorcmd  string
	Shell      string
	Extras     []string
}

func NewConfig(path string, editorMode string, editorCmd string, shell string) (Config, error) {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return Config{}, homeErr
	}

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	}

	return Config{
		Path:       path,
		Editormode: editorMode,
		Editorcmd:  editorCmd,
		Shell:      shell,
	}, nil
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

func ConfigFileMissing(appFs afero.Fs) (bool, error) {
	configFilePath, configFilePathErr := ConfigFilePath()
	if configFilePathErr != nil {
		return false, configFilePathErr
	}

	return FileMissing(appFs, configFilePath), nil
}

func ReadConfigFile(appFs afero.Fs) (Config, error) {
	configFilePath, configFilePathErr := ConfigFilePath()
	if configFilePathErr != nil {
		return Config{}, configFilePathErr
	}

	jsonConfig, readConfigFileErr := ReadFile(appFs, configFilePath)
	if readConfigFileErr != nil {
		return Config{}, readConfigFileErr
	}

	c, fromJsonErr := FromJson(jsonConfig)
	if fromJsonErr != nil {
		return Config{}, fromJsonErr
	}

	return c, nil
}

func FromJson(j []byte) (Config, error) {
	var c Config
	unmarshalErr := json.Unmarshal(j, &c)

	return c, unmarshalErr
}

func SaveExtraDir(appFs afero.Fs, path string) error {
	c, readConfigErr := ReadConfigFile(appFs)
	if readConfigErr != nil {
		return readConfigErr
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
