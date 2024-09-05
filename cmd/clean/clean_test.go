package clean

import (
	"bytes"
	"testing"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
)

func TestCleanCommand(t *testing.T) {
	appFs := afero.NewMemMapFs()

	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		t.Error(aliasFilePathErr)
	}

	t.Run("when alias file is present", func(t *testing.T) {
		var output bytes.Buffer
		pkg.WriteFixtureFile(t, appFs, aliasFilePath, []byte("alias foo=\"cd /bar/baz\""))

		got := RunCleanCmd(appFs, &output)

		if got != nil {
			t.Fatalf("expected 'nil' from RunCleanCmd. got=%q", got)
		}

		if !pkg.FileMissing(appFs, aliasFilePath) {
			t.Fatal("expected alias file to be missing")
		}

		if output.String() != "Aliases file has been deleted." {
			t.Errorf("expected output to say 'Aliases file has been deleted.' got='%s'", output.String())
		}
	})

	t.Run("when alias file is missing", func(t *testing.T) {
		var output bytes.Buffer

		got := RunCleanCmd(appFs, &output)

		if got != nil {
			t.Fatalf("expected 'nil' from RunCleanCmd. got=%q", got)
		}

		if output.String() != "Aliases file is missing." {
			t.Fatalf("expected output to say 'Aliases file is missing.' got='%s'", output.String())
		}
	})
}
