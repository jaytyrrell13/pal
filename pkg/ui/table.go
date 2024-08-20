package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func Table(headers []string, rows [][]string) *table.Table {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("245"))).
		Headers(headers...).
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

	return t
}
