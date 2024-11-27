package pkg

import (
	"testing"

	"github.com/spf13/afero"
)

func TestAliasFilePath(t *testing.T) {
	configDir, configDirErr := ConfigDirPath()
	if configDirErr != nil {
		t.Error(configDirErr)
	}

	got, err := AliasFilePath()
	expected := configDir + "/aliases"

	if got != expected || err != nil {
		t.Errorf("Expected '%q', but got '%q'", expected, got)
	}
}

func TestSaveAliases(t *testing.T) {
	aliases := []Alias{
		{alias: "foo", path: "/foo"},
		{alias: "bar", path: "/bar"},
	}

	t.Run("when aliases is empty", func(t *testing.T) {
		appFs := afero.NewMemMapFs()
		got := SaveAliases(appFs, []Alias{})

		if got != nil {
			t.Errorf("expected 'nil' but got='%q'", got)
		}

		aliasFilePath, aliasFilePathErr := AliasFilePath()
		if aliasFilePathErr != nil {
			t.Errorf("expected 'nil' but got='%q'", aliasFilePathErr)
		}

		if !FileMissing(appFs, aliasFilePath) {
			t.Errorf("expected alias file to be missing")
		}
	})

	t.Run("when aliases is not empty", func(t *testing.T) {
		appFs := afero.NewMemMapFs()

		got := SaveAliases(appFs, aliases)

		if got != nil {
			t.Errorf("expected 'nil' but got='%q'", got)
		}

		aliasFilePath, aliasFilePathErr := AliasFilePath()
		if aliasFilePathErr != nil {
			t.Errorf("expected 'nil' but got='%q'", aliasFilePathErr)
		}

		if FileMissing(appFs, aliasFilePath) {
			t.Error("expected alias file to not be missing")
		}
	})
}

func TestString(t *testing.T) {
	configWithEditorcmd, _ := NewConfig("/foo", "same", "nvim", "zsh")
	configWithoutEditorcmd, _ := NewConfig("/foo", "same", "", "zsh")

	cases := []struct {
		name     string
		config   Config
		expected string
	}{
		{"when Editorcmd is nvim", configWithEditorcmd, "alias foo=\"cd /foo\"\nalias efoo=\"cd /foo && nvim\"\n"},
		{"when Editorcmd is blank", configWithoutEditorcmd, "alias foo=\"cd /foo\"\n"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAlias("foo", "/foo", tt.config.Editorcmd)

			got := a.String()

			if got != tt.expected {
				t.Errorf("expected '%q' but got '%q'", tt.expected, got)
			}
		})
	}
}
