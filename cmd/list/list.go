package list

import (
	"os"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Display aliases in `.pal`",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		AppFs := afero.NewOsFs()

		aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
		cobra.CheckErr(aliasFilePathErr)

		aliasFile, aliasFileErr := pkg.ReadFile(AppFs, aliasFilePath)
		cobra.CheckErr(aliasFileErr)

		os.Stdout.Write(aliasFile)
	},
}
