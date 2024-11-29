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

		pkg.WriteFixtureFile(t, appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"Editormode\": \"editorMode\", \"Editorcmd\": \"editorCmd\", \"Shell\": \"Zsh\", \"Extras\": [\"/foo/bar\"]}"))

		got := RunConfigListCmd(appFs, &output)

		if got != nil {
			t.Fatalf("expected 'nil' from RunConfigListCmd. got=%q", got)
		}

		outputString := output.String()

		if !strings.Contains(outputString, "/foo") {
			t.Fatalf("expected output to contain '/foo': \n%s", outputString)
		}

		if !strings.Contains(outputString, "editorCmd") {
			t.Fatalf("expected output to contain 'editorCmd': \n%s", outputString)
		}

		if !strings.Contains(outputString, "editorMode") {
			t.Fatalf("expected output to contain 'editorMode': \n%s", outputString)
		}

		if !strings.Contains(outputString, "Zsh") {
			t.Fatalf("expected output to contain 'Zsh': \n%s", outputString)
		}

		if !strings.Contains(outputString, "/foo/bar") {
			t.Fatalf("expected output to contain '/foo/bar: \n%s", outputString)
		}
	})
}
