package list

import (
	"os"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Display aliases in `~/pal`",
	Aliases: []string{"ls"},
	RunE: func(_ *cobra.Command, _ []string) error {
		return RunListCmd()
	},
}

func RunListCmd() error {
	AppFs := afero.NewOsFs()

	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	aliasFile, aliasFileErr := pkg.ReadFile(AppFs, aliasFilePath)
	if aliasFileErr != nil {
		return aliasFileErr
	}

	os.Stdout.Write(aliasFile)

	return nil
}
