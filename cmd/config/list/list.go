package list

import (
	"io"
	"os"
	"strings"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var ConfigListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List pal configs",
	Aliases: []string{"ls"},
	RunE: func(_ *cobra.Command, _ []string) error {
		appFs := afero.NewOsFs()

		return RunConfigListCmd(appFs, os.Stdout)
	},
}

func RunConfigListCmd(appFs afero.Fs, w io.Writer) error {
	c, readConfigErr := pkg.ReadConfigFile(appFs)
	if readConfigErr != nil {
		return readConfigErr
	}

	headers := []string{"Key", "Value"}
	rows := [][]string{
		{"Path", c.Path},
		{"Editormode", c.Editormode},
		{"Editorcmd", c.Editorcmd},
		{"Shell", c.Shell},
		{"Extras", strings.Join(c.Extras, ", ")},
	}

	t := ui.Table(headers, rows)

	_, writeErr := w.Write([]byte(t.String()))
	if writeErr != nil {
		return writeErr
	}

	return nil
}
