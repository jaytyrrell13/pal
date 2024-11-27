package add

import (
	"fmt"

	"github.com/jaytyrrell13/pal/cmd/make"
	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	nameFlag      string
	pathFlag      string
	editorCmdFlag string
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an alias for an additional directory",
	RunE: func(_ *cobra.Command, _ []string) error {
		appFs := afero.NewOsFs()

		return RunAddCmd(appFs, nameFlag, pathFlag, editorCmdFlag)
	},
}

func init() {
	AddCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "Name of the additional alias")
	AddCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to your additional directory")
	AddCmd.Flags().StringVarP(&editorCmdFlag, "editorCmd", "e", "", "Editor command e.g. (nvim, subl, code)")
}

func RunAddCmd(appFs afero.Fs, name string, path string, editorCmd string) error {
	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	c, readConfigErr := pkg.ReadConfigFile(appFs)
	if readConfigErr != nil {
		return readConfigErr
	}

	if pkg.FileMissing(appFs, aliasFilePath) {
		runMake, confirmErr := ui.Confirm("Alias file is missing. Would you like to run make command now?")

		if confirmErr != nil {
			return confirmErr
		}

		if !runMake {
			fmt.Println("Please run `pal make` command manually.")
			return nil
		}

		fmt.Println("Running make command.")

		makeCmdErr := make.RunMakeCmd(appFs)
		if makeCmdErr != nil {
			return makeCmdErr
		}
	}

	if name == "" {
		nameString, nameErr := ui.Input("What is the name of the alias?", "foo")

		if nameErr != nil {
			return nameErr
		}

		name = nameString
	}

	if path == "" {
		pathString, pathErr := ui.Input("What is the path for the alias?", "/Users/john/Documents")

		if pathErr != nil {
			return pathErr
		}

		path = pathString
	}

	if editorCmd == "" {
		if c.Editormode == "same" {
			editorCmd = c.Editorcmd
		}

		if c.Editormode == "unique" {
			editorCmdString, editorCmdErr := ui.Input(fmt.Sprintf("What is the editor command for (%s)?", name), "nvim, subl, code")

			if editorCmdErr != nil {
				return editorCmdErr
			}

			editorCmd = editorCmdString
		}
	}

	saveExtraDirErr := pkg.SaveExtraDir(appFs, path)
	if saveExtraDirErr != nil {
		return saveExtraDirErr
	}

	a := pkg.NewAlias(name, path, editorCmd)

	return pkg.AppendToFile(appFs, aliasFilePath, []byte(a.String()))
}
