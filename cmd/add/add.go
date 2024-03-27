package add

import (
	"os"

	"github.com/jaytyrrell13/pal/pkg/aliases"
	"github.com/jaytyrrell13/pal/pkg/config"
	"github.com/jaytyrrell13/pal/pkg/prompts"
	"github.com/spf13/cobra"
)

var (
	nameFlag string
	pathFlag string
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an alias for an additional directory",
	Run: func(cmd *cobra.Command, args []string) {
		name := nameFlag
		path := pathFlag

		if name == "" {
			name = prompts.StringPrompt("What is the name of the alias?")
		}

		if path == "" {
			path = prompts.StringPrompt("What is the path for the alias?")
		}

		config.SaveExtraDir(path)

		if aliases.AliasFileMissing() {
			return
		}

		aliasesFile, openAliasesFileErr := os.OpenFile(aliases.AliasFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
		cobra.CheckErr(openAliasesFileErr)

		c := config.ReadConfigFile()

		var output string
		output += aliases.MakeAliasCommands(name, path, c)

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
	AddCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "Name of the additional alias")
	AddCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to your additional directory")
}
