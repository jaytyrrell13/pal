package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jaytyrrell13/pal/internal"
	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/jaytyrrell13/pal/internal/messages"
	"github.com/spf13/afero"
)

func TestCheckListPrerequisites(t *testing.T) {
	t.Run("when config file does not exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		err := CheckListPrerequisites(fs)

		if err == nil {
			t.Error("expected an error but received 'nil'")
		}
	})

	t.Run("when config file exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		internal.WriteConfigFile(t, fs, config.NewConfig("zsh"))

		err := CheckListPrerequisites(fs)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}
	})
}

func TestRunListCmd(t *testing.T) {
	t.Run("when config does not contain aliases", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		c := config.Config{
			Shell:   "zsh",
			Aliases: []alias.Alias{},
		}
		internal.WriteConfigFile(t, fs, c)

		var output bytes.Buffer
		err := RunListCmd(fs, &output)
		if err == nil {
			t.Errorf("expected '%s' but got %s", messages.Errors["aliasesEmpty"], err)
		}
	})

	t.Run("when config contains aliases", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		c := config.Config{
			Shell: "zsh",
			Aliases: []alias.Alias{
				{Name: "foo", Command: "cd /some/thing"},
				{Name: "bar", Command: "cd /another/thing"},
			},
		}
		internal.WriteConfigFile(t, fs, c)

		var output bytes.Buffer
		err := RunListCmd(fs, &output)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}

		if !strings.Contains(output.String(), "foo") {
			t.Errorf("expected output to contain 'foo': \n%s", output.String())
		}

		if !strings.Contains(output.String(), "cd /some/thing") {
			t.Errorf("expected output to contain 'cd /some/thing': \n%s", output.String())
		}

		if !strings.Contains(output.String(), "bar") {
			t.Errorf("expected output to contain 'bar': \n%s", output.String())
		}

		if !strings.Contains(output.String(), "cd /another/thing") {
			t.Errorf("expected output to contain 'cd /another/thing': \n%s", output.String())
		}
	})
}
