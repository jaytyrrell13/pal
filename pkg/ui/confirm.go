package ui

import "github.com/charmbracelet/huh"

func Confirm(title string) (bool, error) {
	var value bool

	err := huh.NewConfirm().
		Title(title).
		Value(&value).
		Affirmative("Yes").
		Negative("No").
		WithTheme(huh.ThemeBase()).
		Run()

	return value, err
}
