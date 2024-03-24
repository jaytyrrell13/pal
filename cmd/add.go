package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jaytyrrell13/pal/pkg/config"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an alias for an additional directory",
	Run: func(cmd *cobra.Command, args []string) {
		name := prompts.StringPrompt("What is the name of the alias?")

		path := prompts.StringPrompt("What is the path for the alias?")

		config.SaveExtraDir(path)

		homedir, homedirErr := os.UserHomeDir()
		cobra.CheckErr(homedirErr)

		_, aliasFileError := os.Stat(homedir + "/.pal")
		if errors.Is(aliasFileError, os.ErrNotExist) {
			return
		}

		c := config.ReadConfigFile()

		aliasesFile, openAliasesFileErr := os.OpenFile(homedir+"/.pal", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
		cobra.CheckErr(openAliasesFileErr)

		var output string
		output += fmt.Sprintf("alias %s=\"cd %s\"\n", name, path)

		if c.EditorCmd != "" {
			output += fmt.Sprintf("alias %s=\"cd %s && %s\"\n", "e"+name, path, c.EditorCmd)
		}

		if _, err := aliasesFile.Write([]byte(output)); err != nil {
			aliasesFile.Close()
			cobra.CheckErr(err)
		}
		if err := aliasesFile.Close(); err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
