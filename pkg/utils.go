package pkg

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

func ReadFile(path string) []byte {
	file, openErr := os.ReadFile(path)
	cobra.CheckErr(openErr)

	return file
}

func FileMissing(path string) bool {
	_, e := os.Stat(path)
	return errors.Is(e, os.ErrNotExist)
}
