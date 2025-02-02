package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/jaytyrrell13/pal/internal"
	"github.com/spf13/afero"
)

func TestRunInstallCmd(t *testing.T) {
	configFilePath, configFilePathErr := internal.ConfigFilePath()
	if configFilePathErr != nil {
		t.Error(configFilePathErr)
	}

	t.Run("when config file does not exist", func(t *testing.T) {
		ip := InstallPrompts{
			shell: "zsh",
		}

		fs := afero.NewMemMapFs()
		err := RunInstallCmd(fs, ip)
		if err != nil {
			t.Errorf("expected 'nil' but got '%s'", err)
		}

		_, statErr := fs.Stat(configFilePath)
		if errors.Is(statErr, os.ErrNotExist) {
			t.Errorf("expected 'nil' but got '%s'", statErr)
		}

		b, readFileErr := afero.ReadFile(fs, configFilePath)
		if readFileErr != nil {
			t.Errorf("expected 'nil' but got '%s'", readFileErr)
		}

		var c internal.Config
		jsonErr := json.Unmarshal(b, &c)
		if jsonErr != nil {
			t.Errorf("expected 'nil' but got '%s'", jsonErr)
		}

		if c.Shell != "zsh" {
			t.Errorf("expected 'zsh' but got '%s'", c.Shell)
		}
	})

	t.Run("when config file already exists", func(t *testing.T) {
		configDirPath, configDirPathErr := internal.ConfigDirPath()
		if configDirPathErr != nil {
			t.Error(configDirPathErr)
		}

		fs := afero.NewMemMapFs()
		mkDirErr := fs.MkdirAll(configDirPath, 0o755)
		if mkDirErr != nil {
			t.Error(mkDirErr)
		}

		config := internal.NewConfig("bash")

		j, jsonErr := json.Marshal(config)
		if jsonErr != nil {
			t.Error(jsonErr)
		}

		writeFileErr := afero.WriteFile(fs, configFilePath, j, 0o644)
		if writeFileErr != nil {
			t.Error(writeFileErr)
		}

		ip := InstallPrompts{
			shell: "zsh",
		}

		err := RunInstallCmd(fs, ip)
		if err == nil {
			t.Error("expected an error but got 'nil'")
		}

		b, readFileErr := afero.ReadFile(fs, configFilePath)
		if readFileErr != nil {
			t.Errorf("expected 'nil' but got '%s'", readFileErr)
		}

		var c internal.Config
		jsonUnmarshalErr := json.Unmarshal(b, &c)
		if jsonUnmarshalErr != nil {
			t.Errorf("expected 'nil' but got '%s'", jsonErr)
		}

		if c.Shell != "bash" {
			t.Errorf("expected 'bash' but got '%s'", c.Shell)
		}
	})
}
