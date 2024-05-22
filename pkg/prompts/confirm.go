package prompts

import "github.com/charmbracelet/huh"

func Confirm(title string) (bool, error) {
	var value bool

	err := huh.NewConfirm().
		Title(title).
		Value(&value).
		Affirmative("Yes").
		Negative("No").
		Run()

	return value, err
}
