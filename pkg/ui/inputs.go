package ui

import (
	"strings"

	"github.com/charmbracelet/huh"
)

func Input(title string, placeholder string) (string, error) {
	var s string

	err := huh.NewInput().
		Title(title).
		Value(&s).
		Placeholder(placeholder).
		Run()

	return strings.TrimSpace(s), err
}
