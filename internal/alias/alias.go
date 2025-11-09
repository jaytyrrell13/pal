package alias

import (
	"fmt"
)

type Alias struct {
	Name    string
	Command string
}

func (a Alias) ToActionCmd() Alias {
	a.Command = "cd " + a.Command

	return a
}

func (a Alias) ToEditCmd() Alias {
	a.Name = "e" + a.Name
	a.Command = "cd " + a.Command + " && nvim"

	return a
}

func (a Alias) String() string {
	return fmt.Sprintf("alias %s=\"%s\"\n", a.Name, a.Command)
}
