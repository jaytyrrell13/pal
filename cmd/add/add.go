package add

import (
	"os"

	"github.com/jaytyrrell13/pal/pkg"
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
		if pkg.AliasFileMissing() {
			cobra.CheckErr("~/.pal file is missing, please run make command first")
		}

		name := nameFlag
		path := pathFlag

		if name == "" {
			name = prompts.StringPrompt("What is the name of the alias?")
		}

		if path == "" {
			path = prompts.StringPrompt("What is the path for the alias?")
		}

		pkg.SaveExtraDir(path)

		aliasesFile, openAliasesFileErr := os.OpenFile(pkg.AliasFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
		cobra.CheckErr(openAliasesFileErr)

		c := pkg.ReadConfigFile()

		var output string
		output += pkg.MakeAliasCommands(name, path, c)

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
