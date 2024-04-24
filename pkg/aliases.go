package pkg

import (
	"fmt"
	"os"
)

func AliasFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()

	return homeDir + "/.pal", err
}

func MakeAliasCommands(name string, path string, config Config) string {
	var output string
	output += fmt.Sprintf("alias %s=\"cd %s\"\n", name, path)

	if config.EditorCmd != "" {
		output += fmt.Sprintf("alias %s=\"cd %s && %s\"\n", "e"+name, path, config.EditorCmd)
	}

	return output
}
