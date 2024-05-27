package clean

import (
	"fmt"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete the pal aliases file",
	RunE: func(_ *cobra.Command, _ []string) error {
		return RunCleanCmd()
	},
}

func RunCleanCmd() error {
	AppFs := afero.NewOsFs()

	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	if pkg.FileMissing(AppFs, aliasFilePath) {
		fmt.Println("Aliases file is missing.")
		return nil
	}

	removeFileErr := pkg.RemoveFile(AppFs, aliasFilePath)
	if removeFileErr != nil {
		return removeFileErr
	}

	fmt.Println("Aliases file has been deleted.")

	return nil
}
