package clean

import (
	"io"
	"os"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete the pal aliases file",
	RunE: func(_ *cobra.Command, _ []string) error {
		appFs := afero.NewOsFs()

		return RunCleanCmd(appFs, os.Stdout)
	},
}

func RunCleanCmd(appFs afero.Fs, w io.Writer) error {
	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	if pkg.FileMissing(appFs, aliasFilePath) {
		_, writeErr := w.Write([]byte("Aliases file is missing."))
		if writeErr != nil {
			return writeErr
		}

		return nil
	}

	removeFileErr := pkg.RemoveFile(appFs, aliasFilePath)
	if removeFileErr != nil {
		return removeFileErr
	}

	_, writeErr := w.Write([]byte("Aliases file been deleted."))
	if writeErr != nil {
		return writeErr
	}

	return nil
}
