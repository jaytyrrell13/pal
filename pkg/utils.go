package pkg

import (
	"errors"
	"io/fs"
	"os"

	"github.com/spf13/afero"
)

func ReadFile(afs afero.Fs, path string) ([]byte, error) {
	return afero.ReadFile(afs, path)
}

func WriteFile(afs afero.Fs, fileName string, data []byte, perm fs.FileMode) error {
	return afero.WriteFile(afs, fileName, data, perm)
}

func RemoveFile(afs afero.Fs, path string) error {
	return afs.Remove(path)
}

func FileMissing(fs afero.Fs, path string) bool {
	_, e := fs.Stat(path)
	return errors.Is(e, os.ErrNotExist)
}
