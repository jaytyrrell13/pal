package cmd

import (
	"testing"

	"github.com/jaytyrrell13/pal/internal"
	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/spf13/afero"
)

func TestCheckUpdatePrerequisites(t *testing.T) {
	t.Run("when config file does not exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		err := CheckUpdatePrerequisites(fs)

		if err == nil {
			t.Error("expected an error but received 'nil'")
		}
	})

	t.Run("when aliases file does not exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		internal.WriteConfigFile(t, fs, config.NewConfig("zsh"))

		err := CheckUpdatePrerequisites(fs)

		if err == nil {
			t.Error("expected an error but received 'nil'")
		}
	})

	t.Run("when config and aliases files exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		internal.WriteConfigFile(t, fs, config.NewConfig("zsh"))
		internal.WriteAliasesFile(t, fs)

		err := CheckUpdatePrerequisites(fs)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}
	})
}

func TestRunUpdateCmd(t *testing.T) {
	t.Run("when given a list of updated aliases", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		c := config.Config{
			Shell: "zsh",
			Aliases: []alias.Alias{
				{Name: "docs", Command: "cd /foo/Documents"},
			},
		}

		internal.WriteConfigFile(t, fs, c)
		internal.WriteAliasesFile(t, fs)

		up := UpdatePrompts{
			updatedAliases: []UpdatedAlias{
				{oldName: "docs", newAlias: alias.Alias{
					Name:    "docz",
					Command: "cd /bar/Documents",
				}},
			},
		}

		err := RunUpdateCmd(fs, up)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}

		c, configErr := config.ReadConfigFile(fs)
		if configErr != nil {
			t.Errorf("expected 'nil' but got error: %s", configErr)
		}

		internal.AssertEquals(t, "docz", c.Aliases[0].Name)
		internal.AssertEquals(t, "cd /bar/Documents", c.Aliases[0].Command)

		internal.AssertAliasFileContains(t, fs, "docz")
		internal.AssertAliasFileContains(t, fs, "cd /bar/Documents")
	})
}
