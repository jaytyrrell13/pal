package cmd

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/spf13/afero"
)

func TestRunCreateCmd(t *testing.T) {
	t.Run("when category is 'parent'", func(t *testing.T) {
		cp := CreatePrompts{
			category: "parent",
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

		mkDirErr := fs.MkdirAll("/foo/Code", 0o755)
		if mkDirErr != nil {
			t.Error(mkDirErr)
		}

		mkDir2Err := fs.Mkdir("/foo/Code/project-one", 0o755)
		if mkDir2Err != nil {
			t.Error(mkDir2Err)
		}

		mkDir3Err := fs.Mkdir("/foo/Code/project-two", 0o755)
		if mkDir3Err != nil {
			t.Error(mkDir3Err)
		}

		err := RunCreateCmd(fs, cp)
		if err != nil {
			t.Errorf("expected 'nil' but got '%s'", err)
		}

		b, readFileErr := afero.ReadFile(fs, configFilePath)
		if readFileErr != nil {
			t.Errorf("expected 'nil' but got '%s'", readFileErr)
		}

		var c config.Config
		jsonUnmarshalErr := json.Unmarshal(b, &c)
		if jsonUnmarshalErr != nil {
			t.Errorf("expected 'nil' but got '%s'", jsonUnmarshalErr)
		}

		if c.Aliases[0].Name != "po" {
			t.Errorf("expected 'po' but got '%s'", c.Aliases[0].Name)
		}

		if c.Aliases[0].Command != "/foo/Code/project-one" {
			t.Errorf("expected '/foo/Code/project-one' but got '%s'", c.Aliases[0].Command)
		}

		if c.Aliases[1].Name != "pt" {
			t.Errorf("expected 'pt' but got '%s'", c.Aliases[0].Name)
		}

		if c.Aliases[1].Command != "/foo/Code/project-two" {
			t.Errorf("expected '/foo/Code/project-two' but got '%s'", c.Aliases[0].Command)
		}
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

		mkDirErr := fs.MkdirAll("/foo/Documents/work/notes", 0o755)
		if mkDirErr != nil {
			t.Error(mkDirErr)
		}

		err := RunCreateCmd(fs, cp)
		if err != nil {
			t.Errorf("expected 'nil' but got '%s'", err)
		}

		b, readFileErr := afero.ReadFile(fs, configFilePath)
		if readFileErr != nil {
			t.Errorf("expected 'nil' but got '%s'", readFileErr)
		}

		var c config.Config
		jsonUnmarshalErr := json.Unmarshal(b, &c)
		if jsonUnmarshalErr != nil {
			t.Errorf("expected 'nil' but got '%s'", jsonUnmarshalErr)
		}

		fmt.Println(c)

		if c.Aliases[0].Name != "wn" {
			t.Errorf("expected 'wn' but got '%s'", c.Aliases[0].Name)
		}

		if c.Aliases[0].Command != "/foo/Documents/work/notes" {
			t.Errorf("expected '/foo/Documents/work/notes' but got '%s'", c.Aliases[0].Command)
		}
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

		err := RunCreateCmd(fs, cp)
		if err != nil {
			t.Errorf("expected 'nil' but got '%s'", err)
		}

		b, readFileErr := afero.ReadFile(fs, configFilePath)
		if readFileErr != nil {
			t.Errorf("expected 'nil' but got '%s'", readFileErr)
		}

		var c config.Config
		jsonUnmarshalErr := json.Unmarshal(b, &c)
		if jsonUnmarshalErr != nil {
			t.Errorf("expected 'nil' but got '%s'", jsonUnmarshalErr)
		}

		fmt.Println(c)

		if c.Aliases[0].Name != "ll" {
			t.Errorf("expected 'll' but got '%s'", c.Aliases[0].Name)
		}

		if c.Aliases[0].Command != "ls -lah" {
			t.Errorf("expected 'ls -lah' but got '%s'", c.Aliases[0].Command)
		}
	})
}
