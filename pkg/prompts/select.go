package prompts

import "github.com/charmbracelet/huh"

func Select(title string) (string, error) {
	var s string

	err := huh.NewSelect[string]().
		Title(title).
		Options(huh.NewOptions("Bash/ZSH", "Fish")...).
		Value(&s).
		Run()

	return s, err
}
