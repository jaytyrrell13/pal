package aliases

import (
	"os"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/cobra"
)

var AliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "Display aliases in `.pal`",
	Run: func(cmd *cobra.Command, args []string) {
		os.Stdout.Write(pkg.ReadAliasFile())
	},
}
