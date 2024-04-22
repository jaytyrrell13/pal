package pkg

import (
	"errors"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func ReadFile(afs afero.Fs, path string) []byte {
	file, openErr := afero.ReadFile(afs, path)
	cobra.CheckErr(openErr)

	return file
}

func FileMissing(fs afero.Fs, path string) bool {
	_, e := fs.Stat(path)
	return errors.Is(e, os.ErrNotExist)
}
