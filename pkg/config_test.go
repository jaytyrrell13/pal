package pkg

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestConfigDirPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()

	got := ConfigDirPath()
	expected := homeDir + "/.config/pal"

	if got != expected || err != nil {
		t.Errorf("Expected '%q', but got '%q'", expected, got)
	}
}

func TestConfigFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()

	got := ConfigFilePath()
	expected := homeDir + "/.config/pal/config.json"

	if got != expected || err != nil {
		t.Errorf("Expected '%q', but got '%q'", expected, got)
	}
}

func TestMakeConfigDir(t *testing.T) {
	appFs := afero.NewMemMapFs()

	got := MakeConfigDir(appFs)

	if got != nil {
		t.Errorf("Expected 'nil', but got '%q'", got)
	}
}

func TestFromJson(t *testing.T) {
	json := "{\"Path\": \"/foo/bar\", \"EditorCmd\": \"foo\", \"Extras\":null}"
	got, err := FromJson([]byte(json))

	expected := Config{
		Path:      "/foo/bar",
		EditorCmd: "foo",
	}

	if got.Path != expected.Path || got.EditorCmd != expected.EditorCmd || err != nil {
		t.Errorf("Expected Path '%q' EditorCmd '%q', but got Path '%q' EditorCmd '%q'", expected.EditorCmd, expected.Path, got.Path, got.EditorCmd)
	}
}

func TestAsJson(t *testing.T) {
	config := Config{
		Path:      "/foo/bar",
		EditorCmd: "foo",
	}

	got, err := config.AsJson()

	if got == nil || err != nil {
		t.Errorf("Got '%q' Err '%q'", got, err)
	}
}
