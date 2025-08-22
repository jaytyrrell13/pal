package messages

var Errors = map[string]string{
	"configExists":   "config file already exists",
	"configMissing":  "config file does not exist. run `pal install` command to create a new config file",
	"aliasesMissing": "aliases file does not exist. run `pal create` command to begin creating aliases",
	"aliasesEmpty":   "there are no aliases set up. run `pal create` command to begin creating aliases",
}
