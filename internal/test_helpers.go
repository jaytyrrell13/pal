package internal

import (
	"encoding/json"
	"strings"
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

func AssertAliasFileContains(t *testing.T, fs afero.Fs, expected string) {
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

	if !strings.Contains(str, expected) {
		t.Errorf("expected '%s' to contain '%s'", str, expected)
	}
}

func AssertAliasFileDoesntContain(t *testing.T, fs afero.Fs, expected string) {
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

	if strings.Contains(str, expected) {
		t.Errorf("expected '%s' to not contain '%s'", str, expected)
	}
}
