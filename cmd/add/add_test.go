package add

import (
	"strings"
	"testing"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
)

func TestAddCommand(t *testing.T) {
	appFs := afero.NewMemMapFs()

	aliasFilePath, aliasFilePathErr := pkg.AliasFilePath()
	if aliasFilePathErr != nil {
		t.Error(aliasFilePathErr)
	}

	writeFileErr := afero.WriteFile(appFs, aliasFilePath, []byte("alias foo=\"cd /bar/baz\"\n"), 0o755)
	if writeFileErr != nil {
		t.Fatalf("WriteFile Error: '%q'", writeFileErr)
	}

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		t.Error(configFilePathErr)
	}

	writeConfigFileErr := afero.WriteFile(appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"EditorCmd\": \"nvim\"}"), 0o755)
	if writeConfigFileErr != nil {
		t.Fatalf("WriteConfigFile Error: '%q'", writeConfigFileErr)
	}

	got := RunAddCmd(appFs, "bark", "/foo/baz")

	if got != nil {
		t.Errorf("expected 'nil' from RunAddCmd. got=%q", got)
	}

	assertFileContains(t, appFs, aliasFilePath, "bark")
	assertFileContains(t, appFs, aliasFilePath, "/foo/baz")
	assertFileContains(t, appFs, configFilePath, "/foo/baz")
}

func assertFileContains(t *testing.T, appFs afero.Fs, filePath, s string) {
	t.Helper()
	readFile, readFileErr := afero.ReadFile(appFs, filePath)
	if readFileErr != nil {
		t.Errorf("ReadFile Error: %+v", readFileErr)
	}

	if !strings.Contains(string(readFile), s) {
		t.Errorf("expected '%s' to be in file '%s'", s, filePath)
	}
}
