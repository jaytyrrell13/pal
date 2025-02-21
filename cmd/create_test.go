package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/spf13/afero"
)

func TestCheckCreatePrerequisites(t *testing.T) {
	t.Run("when config file does not exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		err := CheckCreatePrerequisites(fs)

		if err == nil {
			t.Error("expected an error but received 'nil'")
		}
	})

	t.Run("when config file does exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		writeConfigFile(t, fs)

		err := CheckCreatePrerequisites(fs)
		if err != nil {
			t.Errorf("expected 'nil' but got error: %s", err)
		}
	})
}

func TestRunCreateCmd(t *testing.T) {
	t.Run("when category is 'parent'", func(t *testing.T) {
		cp := CreatePrompts{
			category: "parent",
			editCmd:  "yes",
			aliases: []alias.Alias{
				{
					Name:    "po",
					Command: "/foo/Code/project-one",
				},
				{
					Name:    "pt",
					Command: "/foo/Code/project-two",
				},
			},
		}

		fs := afero.NewMemMapFs()

		writeConfigFile(t, fs)

		makeDirAll(t, fs, "/foo/Code/project-one")
		makeDirAll(t, fs, "/foo/Code/project-two")

		err := RunCreateCmd(fs, cp)
		if err != nil {
			t.Errorf("expected 'nil' but got '%s'", err)
		}

		c, configErr := config.ReadConfigFile(fs)
		if configErr != nil {
			t.Errorf("expected 'nil' but got '%s'", configErr)
		}

		assertAliasMatches(t, "po", c.Aliases[0].Name)
		assertAliasMatches(t, "cd /foo/Code/project-one", c.Aliases[0].Command)
		assertAliasMatches(t, "epo", c.Aliases[1].Name)
		assertAliasMatches(t, "cd /foo/Code/project-one && nvim", c.Aliases[1].Command)
		assertAliasMatches(t, "pt", c.Aliases[2].Name)
		assertAliasMatches(t, "cd /foo/Code/project-two", c.Aliases[2].Command)
		assertAliasMatches(t, "ept", c.Aliases[3].Name)
		assertAliasMatches(t, "cd /foo/Code/project-two && nvim", c.Aliases[3].Command)

		assertAliasFileContains(t, fs, "po")
		assertAliasFileContains(t, fs, "cd /foo/Code/project-one")
		assertAliasFileContains(t, fs, "epo")
		assertAliasFileContains(t, fs, "cd /foo/Code/project-one && nvim")
		assertAliasFileContains(t, fs, "pt")
		assertAliasFileContains(t, fs, "cd /foo/Code/project-two")
		assertAliasFileContains(t, fs, "ept")
		assertAliasFileContains(t, fs, "cd /foo/Code/project-two && nvim")
	})

	t.Run("when category is 'directory'", func(t *testing.T) {
		cp := CreatePrompts{
			category: "directory",
			aliases: []alias.Alias{
				{
					Name:    "wn",
					Command: "/foo/Documents/work/notes",
				},
			},
		}

		fs := afero.NewMemMapFs()

		writeConfigFile(t, fs)

		makeDirAll(t, fs, "/foo/Documents/work/notes")

		err := RunCreateCmd(fs, cp)
		if err != nil {
			t.Errorf("expected 'nil' but got '%s'", err)
		}

		c, configErr := config.ReadConfigFile(fs)
		if configErr != nil {
			t.Errorf("expected 'nil' but got '%s'", configErr)
		}

		assertAliasMatches(t, "wn", c.Aliases[0].Name)
		assertAliasMatches(t, "cd /foo/Documents/work/notes", c.Aliases[0].Command)

		assertAliasFileContains(t, fs, "wn")
		assertAliasFileContains(t, fs, "cd /foo/Documents/work/notes")
	})

	t.Run("when category is 'action'", func(t *testing.T) {
		cp := CreatePrompts{
			category: "action",
			aliases: []alias.Alias{
				{
					Name:    "ll",
					Command: "ls -lah",
				},
			},
		}

		fs := afero.NewMemMapFs()

		writeConfigFile(t, fs)

		err := RunCreateCmd(fs, cp)
		if err != nil {
			t.Errorf("expected 'nil' but got '%s'", err)
		}

		c, configErr := config.ReadConfigFile(fs)
		if configErr != nil {
			t.Errorf("expected 'nil' but got '%s'", configErr)
		}

		assertAliasMatches(t, "ll", c.Aliases[0].Name)
		assertAliasMatches(t, "ls -lah", c.Aliases[0].Command)

		assertAliasFileContains(t, fs, "ll")
		assertAliasFileContains(t, fs, "ls -lah")
	})

	t.Run("when config does not exist", func(t *testing.T) {
		cp := CreatePrompts{
			category: "action",
			aliases: []alias.Alias{
				{
					Name:    "ll",
					Command: "ls -lah",
				},
			},
		}

		fs := afero.NewMemMapFs()

		err := RunCreateCmd(fs, cp)
		if err == nil {
			t.Error("expected an error but got 'nil'")
		}
	})
}

func assertAliasMatches(t *testing.T, expected string, actual string) {
	t.Helper()

	if actual != expected {
		t.Errorf("expected '%s' but got '%s'", expected, actual)
	}
}

func writeConfigFile(t *testing.T, fs afero.Fs) {
	t.Helper()

	configFilePath, configFilePathErr := config.ConfigFilePath()
	if configFilePathErr != nil {
		t.Error(configFilePathErr)
	}

	j, jsonErr := json.Marshal(config.NewConfig("bash"))
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	writeFileErr := afero.WriteFile(fs, configFilePath, j, 0o644)
	if writeFileErr != nil {
		t.Error(writeFileErr)
	}
}

func makeDirAll(t *testing.T, fs afero.Fs, path string) {
	t.Helper()

	mkDirErr := fs.MkdirAll(path, 0o755)
	if mkDirErr != nil {
		t.Error(mkDirErr)
	}
}

func assertAliasFileContains(t *testing.T, fs afero.Fs, expected string) {
	t.Helper()

	aliasFilePath, aliasFilePathErr := config.AliasFilePath()
	if aliasFilePathErr != nil {
		t.Error(aliasFilePathErr)
	}
	b, readFileErr := afero.ReadFile(fs, aliasFilePath)
	if readFileErr != nil {
		t.Error(readFileErr)
	}

	str := string(b)
	strings.Contains(str, expected)
}
