package clean

import (
	"fmt"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete `/.pal` file in your home directory",
	Run: func(cmd *cobra.Command, args []string) {
		err := RunCleanCmd()
		cobra.CheckErr(err)
	},
}

func RunCleanCmd() error {
	AppFs := afero.NewOsFs()

	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	if pkg.FileMissing(AppFs, aliasFilePath) {
		fmt.Println("~/.pal file is missing.")
		return nil
	}

	removeFileErr := pkg.RemoveFile(AppFs, aliasFilePath)
	if removeFileErr != nil {
		return removeFileErr
	}

	fmt.Println("~/.pal file has been deleted.")

	return nil
}
