package config

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jaytyrrell13/pal/internal"
	cfg "github.com/jaytyrrell13/pal/internal/config"
	"github.com/spf13/afero"
)

func TestCheckListPrerequisites(t *testing.T) {
	t.Run("when config file does not exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		err := CheckListPrerequisites(fs)

		if err == nil {
			t.Error("expected an error but received 'nil'")
		}
	})

	t.Run("when config file does exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		internal.WriteConfigFile(t, fs, cfg.NewConfig("zsh"))

		err := CheckListPrerequisites(fs)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}
	})
}

func TestRunListCmd(t *testing.T) {
	t.Run("lists config file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		var output bytes.Buffer

		internal.WriteConfigFile(t, fs, cfg.NewConfig("zsh"))

		got := RunListCmd(fs, &output)

		if got != nil {
			t.Fatalf("expected 'nil' from RunConfigListCmd. got=%q", got)
		}

		outputString := output.String()

		if !strings.Contains(outputString, "zsh") {
			t.Fatalf("expected output to contain 'zsh': \n%s", outputString)
		}
	})
}
