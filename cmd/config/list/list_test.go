package list

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
)

func TestConfigListCommand(t *testing.T) {
	appFs := afero.NewMemMapFs()
	var output bytes.Buffer

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		t.Fatalf("ConfigDirPath Error: '%q'", configFilePathErr)
	}

	writeFileErr := afero.WriteFile(appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"EditorCmd\": \"editorCmd\"}"), 0o755)
	if writeFileErr != nil {
		t.Fatalf("WriteFile Error: '%q'", writeFileErr)
	}

	got := RunConfigListCmd(appFs, &output)

	if got != nil {
		t.Fatalf("expected 'nil' from RunConfigListCmd. got=%q", got)
	}

	if !strings.Contains(output.String(), "/foo") {
		t.Fatalf("expected output to contain '/foo': \n%s", output.String())
	}

	if !strings.Contains(output.String(), "editorCmd") {
		t.Fatalf("expected output to contain 'editorCmd': \n%s", output.String())
	}
}
