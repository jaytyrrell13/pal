package alias

type Alias struct {
	Name    string
	Command string
}

func (a Alias) ForActionCmd() Alias {
	a.Command = "cd " + a.Command

	return a
}

func (a Alias) ForEditCmd() Alias {
	a.Name = "e" + a.Name
	a.Command = "cd " + a.Command + " && nvim"

	return a
}
