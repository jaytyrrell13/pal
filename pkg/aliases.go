package pkg

import (
	"fmt"

	"github.com/spf13/afero"
)

type Alias struct {
	alias string
	path  string
}

func NewAlias(alias string, path string) Alias {
	return Alias{
		alias: alias,
		path:  path,
	}
}

func AliasFilePath() (string, error) {
	path, err := ConfigDirPath()
	if err != nil {
		return "", err
	}

	return path + "/aliases", err
}

func SaveAliases(appFs afero.Fs, aliases []Alias, c Config) error {
	if len(aliases) == 0 {
		return nil
	}

	aliasFilePath, aliasFilePathErr := AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	var output string
	for _, a := range aliases {
		output += a.String(c)
	}

	return WriteFile(appFs, aliasFilePath, []byte(output), 0o755)
}

func (a *Alias) String(c Config) string {
	var output string
	output += a.asGoToString()

	if c.EditorCmd != "" {
		output += a.asEditString(c.EditorCmd)
	}

	return output
}

func (a *Alias) asGoToString() string {
	return fmt.Sprintf("alias %s=\"cd %s\"\n", a.alias, a.path)
}

func (a *Alias) asEditString(editorCmd string) string {
	return fmt.Sprintf("alias %s=\"cd %s && %s\"\n", "e"+a.alias, a.path, editorCmd)
}
