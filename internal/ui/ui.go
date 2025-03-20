package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func Select(title string, options []huh.Option[string]) (string, error) {
	var value string
	err := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(&value).
		WithTheme(huh.ThemeBase()).
		Run()
	if err != nil {
		return "", err
	}

	return value, nil
}

func MultiSelect(title string, options []huh.Option[string]) ([]string, error) {
	var value []string
	err := huh.NewMultiSelect[string]().
		Title(title).
		Options(options...).
		Value(&value).
		WithTheme(huh.ThemeBase()).
		Run()
	if err != nil {
		return []string{""}, err
	}

	return value, nil
}

type InputProps struct {
	Title string
	Value string
}

func Input(ip InputProps) (string, error) {
	value := ip.Value
	err := huh.NewInput().
		Title(ip.Title).
		Value(&value).
		WithTheme(huh.ThemeBase()).
		Run()
	if err != nil {
		return "", err
	}

	return value, nil
}

func Table(headers []string, rows [][]string) *table.Table {
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("245"))).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 2)
		}).
		Rows(rows...)
}
