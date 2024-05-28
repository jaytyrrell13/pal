package pkg

import (
	"fmt"
)

func AliasFilePath() (string, error) {
	path, err := ConfigDirPath()
	if err != nil {
		return "", err
	}

	return path + "/aliases", err
}

func MakeAliasCommands(name string, path string, config Config) string {
	var output string
	output += fmt.Sprintf("alias %s=\"cd %s\"\n", name, path)

	if config.EditorCmd != "" {
		output += fmt.Sprintf("alias %s=\"cd %s && %s\"\n", "e"+name, path, config.EditorCmd)
	}

	return output
}
