package internal

import (
	"encoding/json"
	"testing"

	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/spf13/afero"
)

func WriteConfigFile(t *testing.T, fs afero.Fs, c config.Config) {
	t.Helper()

	configFilePath, configFilePathErr := config.ConfigFilePath()
	if configFilePathErr != nil {
		t.Error(configFilePathErr)
	}

	j, jsonErr := json.Marshal(c)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	writeFileErr := afero.WriteFile(fs, configFilePath, j, 0o644)
	if writeFileErr != nil {
		t.Error(writeFileErr)
	}
}

func WriteAliasesFile(t *testing.T, fs afero.Fs) {
	t.Helper()

	aliasFilePath, aliasFilePathErr := config.AliasFilePath()
	if aliasFilePathErr != nil {
		t.Error(aliasFilePathErr)
	}

	alias := alias.Alias{
		Name: "docs", Command: "cd /foo/Documents",
	}

	writeFileErr := afero.WriteFile(fs, aliasFilePath, []byte(alias.String()), 0o644)
	if writeFileErr != nil {
		t.Error(writeFileErr)
	}
}
