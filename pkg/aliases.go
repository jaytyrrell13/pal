package pkg

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func AliasFilePath() string {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return homeDir + "/.pal"
}

func AliasFileMissing() bool {
	_, e := os.Stat(AliasFilePath())
	return errors.Is(e, os.ErrNotExist)
}

func MakeAliasCommands(name string, path string, config Config) string {
	var output string
	output += fmt.Sprintf("alias %s=\"cd %s\"\n", name, path)

	if config.EditorCmd != "" {
		output += fmt.Sprintf("alias %s=\"cd %s && %s\"\n", "e"+name, path, config.EditorCmd)
	}

	return output
}
