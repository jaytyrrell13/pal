package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

func StringPrompt(label string, stdin io.Reader) string {
	var s string
	r := bufio.NewReader(stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	return strings.TrimSpace(s)
}

func Input(title string, placeholder string) (string, error) {
	var s string

	err := huh.NewInput().
		Title(title).
		Value(&s).
		Placeholder(placeholder).
		Run()

	return strings.TrimSpace(s), err
}
