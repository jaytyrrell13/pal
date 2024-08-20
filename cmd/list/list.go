package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

	aliases := strings.TrimSpace(string(aliasFile[:]))

	rows := [][]string{}
	splitAliases := strings.Split(aliases, "\n")
	for _, a := range splitAliases {
		trimmed := strings.TrimPrefix(a, "alias ")
		split := strings.Split(trimmed, "=")

		rows = append(rows, []string{split[0], split[1]})
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("245"))).
		Headers("Alias", "Command").
		StyleFunc(func(row, col int) lipgloss.Style {
			baseStyle := lipgloss.NewStyle().Padding(0, 2)

			switch {
			case row == 0:
				return baseStyle.Bold(true)
			case row%2 == 0:
				return baseStyle.Foreground(lipgloss.Color("240"))
			default:
				return baseStyle
			}
		}).
		Rows(rows...)

	fmt.Println(t)

	return nil
}
