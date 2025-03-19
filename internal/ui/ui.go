package ui

import "github.com/charmbracelet/huh"

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
