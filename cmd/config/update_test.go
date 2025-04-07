package config

import (
	"testing"

	"github.com/jaytyrrell13/pal/internal"
	cfg "github.com/jaytyrrell13/pal/internal/config"
	"github.com/spf13/afero"
)

func TestCheckUpdatePrequisites(t *testing.T) {
	t.Run("when config file does not exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		err := CheckUpdatePrerequisites(fs)

		if err == nil {
			t.Error("expected an error but received 'nil'")
		}
	})

	t.Run("when config file exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		internal.WriteConfigFile(t, fs, cfg.NewConfig("zsh"))

		err := CheckUpdatePrerequisites(fs)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}
	})
}

func TestRunUpdateCmd(t *testing.T) {
	t.Run("when given a list of updated config", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		internal.WriteConfigFile(t, fs, cfg.NewConfig("bash"))

		up := UpdatePrompts{
			updatedConfig: []UpdatedConfig{
				{oldKey: "Shell", newValue: "zsh"},
			},
		}

		err := RunUpdateCmd(fs, up)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}

		c, configErr := cfg.ReadConfigFile(fs)
		if configErr != nil {
			t.Errorf("expected 'nil' but got error: %s", configErr)
		}

		internal.AssertEquals(t, "zsh", c.Shell)
	})
}
