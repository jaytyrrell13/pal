package cmd

import (
	"slices"
	"testing"

	"github.com/jaytyrrell13/pal/internal"
	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/spf13/afero"
)

func TestCheckRemovePrerequisites(t *testing.T) {
	t.Run("when config file does not exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		err := CheckRemovePrerequisites(fs)

		if err == nil {
			t.Error("expected an error but received 'nil'")
		}
	})

	t.Run("when aliases file does not exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		internal.WriteConfigFile(t, fs, config.NewConfig("zsh"))

		err := CheckRemovePrerequisites(fs)

		if err == nil {
			t.Error("expected an error but received 'nil'")
		}
	})

	t.Run("when config and aliases files exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		internal.WriteConfigFile(t, fs, config.NewConfig("zsh"))
		internal.WriteAliasesFile(t, fs)

		err := CheckRemovePrerequisites(fs)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}
	})
}

func TestRunRemoveCmd(t *testing.T) {
	t.Run("when given a list of aliases to remove", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		dp := RemovePrompts{
			aliasesToRemove: []string{"bar"},
		}

		c := config.Config{
			Shell: "zsh",
			Aliases: []alias.Alias{
				{Name: "foo", Command: "cd /some/thing"},
				{Name: "bar", Command: "cd /another/thing"},
			},
		}

		internal.WriteConfigFile(t, fs, c)
		internal.WriteAliasesFile(t, fs)

		err := RunRemoveCmd(fs, dp)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}

		c, configErr := config.ReadConfigFile(fs)
		if configErr != nil {
			t.Errorf("expected 'nil' but got error: %s", configErr)
		}

		var aliasNames []string
		for _, a := range c.Aliases {
			aliasNames = append(aliasNames, a.Name)
		}

		if slices.Contains(aliasNames, "bar") {
			t.Error("expected 'bar' to be removed")
		}

		internal.AssertAliasFileContains(t, fs, "foo")
		internal.AssertAliasFileDoesntContain(t, fs, "bar")
	})
}
