package list

import (
	"testing"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
)

func TestListCommand(t *testing.T) {
	appFs := afero.NewMemMapFs()
	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		t.Error(aliasFilePathErr)
	}

	writeFileErr := afero.WriteFile(appFs, aliasFilePath, []byte("alias foo=\"cd /foo/bar\""), 0o755)
	if writeFileErr != nil {
		t.Fatalf("WriteFile Error: '%q'", writeFileErr)
	}

	got := RunListCmd(appFs)

	if got != nil {
		t.Fatalf("expected 'nil' from RunListCmd. got=%q", got)
	}
}
