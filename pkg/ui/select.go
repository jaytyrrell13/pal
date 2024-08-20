package ui

import "github.com/charmbracelet/huh"

func Select(title string, options []huh.Option[string]) (string, error) {
	var s string

	err := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(&s).
		Run()

	return s, err
}
