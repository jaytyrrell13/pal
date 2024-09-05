package pkg

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func ReadFile(afs afero.Fs, path string) ([]byte, error) {
	return afero.ReadFile(afs, path)
}

func AppendToFile(afs afero.Fs, path string, data []byte) error {
	fileContainsBytes, fileContainsBytesErr := afero.FileContainsBytes(afs, path, data)
	if fileContainsBytesErr != nil {
		return fileContainsBytesErr
	}

	if fileContainsBytes {
		return nil
	}

	f, openErr := afs.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
	if openErr != nil {
		return openErr
	}

	if _, writeErr := f.Write(data); writeErr != nil {
		f.Close()
		return writeErr
	}

	return nil
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

func WriteFixtureFile(t *testing.T, appFs afero.Fs, path string, data []byte) {
	t.Helper()
	writeFileErr := afero.WriteFile(appFs, path, data, 0o755)
	if writeFileErr != nil {
		t.Fatalf("WriteFixtureFile Error: '%q'", writeFileErr)
	}
}
