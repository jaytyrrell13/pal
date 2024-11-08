package list

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
)

func TestConfigListCommand(t *testing.T) {
	t.Run("returns an error when config file is missing", func(t *testing.T) {
		appFs := afero.NewMemMapFs()
		var output bytes.Buffer

		got := RunConfigListCmd(appFs, &output)

		if got == nil {
			t.Error("expected an error but got 'nil'")
		}
	})

	t.Run("lists config file", func(t *testing.T) {
		appFs := afero.NewMemMapFs()
		var output bytes.Buffer

		configFilePath, configFilePathErr := pkg.ConfigFilePath()
		if configFilePathErr != nil {
			t.Fatalf("ConfigDirPath Error: '%q'", configFilePathErr)
		}

		pkg.WriteFixtureFile(t, appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"Editorcmd\": \"editorCmd\"}"))

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
	})
}
