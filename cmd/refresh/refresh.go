package refresh

import (
	"fmt"

	"github.com/jaytyrrell13/pal/cmd/make"
	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var RefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Delete pal aliases file and run `make` command",
	RunE: func(_ *cobra.Command, _ []string) error {
		appFs := afero.NewOsFs()

		return RunRefreshCmd(appFs)
	},
}

func RunRefreshCmd(appFs afero.Fs) error {
	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		return aliasFilePathErr
	}

	if pkg.FileMissing(appFs, aliasFilePath) {
		fmt.Println("Aliases file is missing")
		return nil
	}

	fmt.Println("Removing pal aliases file")
	removeFileErr := pkg.RemoveFile(appFs, aliasFilePath)
	if removeFileErr != nil {
		return removeFileErr
	}

	fmt.Println("Running `make` command")
	return make.RunMakeCmd(appFs)
}
