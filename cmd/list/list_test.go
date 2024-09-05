package list

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
)

func TestListCommand(t *testing.T) {
	appFs := afero.NewMemMapFs()
	var output bytes.Buffer

	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		t.Error(aliasFilePathErr)
	}

	pkg.WriteFixtureFile(t, appFs, aliasFilePath, []byte("alias foo=\"cd /bar/baz\""))

	got := RunListCmd(appFs, &output)

	if got != nil {
		t.Fatalf("expected 'nil' from RunListCmd. got=%q", got)
	}

	if !strings.Contains(output.String(), "foo") {
		t.Fatalf("expected output to contain 'foo': \n%s", output.String())
	}

	if !strings.Contains(output.String(), "/bar/baz") {
		t.Fatalf("expected output to contain '/bar/baz': \n%s", output.String())
	}
}
