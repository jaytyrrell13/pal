package list

import (
	"fmt"
	"strings"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/ui"
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

	aliases := strings.TrimSpace(string(aliasFile[:]))

	headers := []string{"Alias", "Command"}
	rows := [][]string{}
	splitAliases := strings.Split(aliases, "\n")
	for _, a := range splitAliases {
		trimmed := strings.TrimPrefix(a, "alias ")
		split := strings.Split(trimmed, "=")

		rows = append(rows, []string{split[0], split[1]})
	}

	t := ui.Table(headers, rows)

	fmt.Println(t)

	return nil
}
