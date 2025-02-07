package ui

import "github.com/charmbracelet/huh"

func Select(title string, options []huh.Option[string]) (string, error) {
	var value string
	err := huh.NewSelect[string]().
		Title("What shell do you use?").
		Options(options...).
		Value(&value).
		WithTheme(huh.ThemeBase()).
		Run()
	if err != nil {
		return "", err
	}

	return value, nil
}

func Input(title string) (string, error) {
	var value string
	err := huh.NewInput().
		Title(title).
		Value(&value).
		WithTheme(huh.ThemeBase()).
		Run()
	if err != nil {
		return "", err
	}

	return value, nil
}
