package pkg

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestConfigDirPath(t *testing.T) {
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		t.Error(homeDirErr)
	}

	got, err := ConfigDirPath()
	expected := homeDir + "/.config/pal"

	if got != expected || err != nil {
		t.Errorf("Expected '%q', but got '%q'", expected, got)
	}
}

func TestConfigFilePath(t *testing.T) {
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		t.Error(homeDirErr)
	}

	got, err := ConfigFilePath()
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

func TestSaveExtraDir(t *testing.T) {
	appFs := afero.NewMemMapFs()

	configFilePath, configFileErr := ConfigFilePath()
	if configFileErr != nil {
		t.Error(configFileErr)
	}

	writeFileErr := afero.WriteFile(appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"EditorCmd\": \"bar\"}"), 0o755)
	if writeFileErr != nil {
		t.Fatalf("WriteFile Error: '%q'", writeFileErr)
	}

	got := SaveExtraDir(appFs, "/bar/baz")

	if got != nil {
		t.Errorf("Expected 'nil', but got '%q'", got)
	}
}

func TestFromJson(t *testing.T) {
	json := "{\"Path\": \"/foo/bar\", \"EditorCmd\": \"foo\", \"Extras\":null}"
	got, err := FromJson([]byte(json))

	expected := NewConfig("/foo/bar", "foo", "")

	if got.Path != expected.Path || got.EditorCmd != expected.EditorCmd || err != nil {
		t.Errorf("Expected Path '%q' EditorCmd '%q', but got Path '%q' EditorCmd '%q'", expected.EditorCmd, expected.Path, got.Path, got.EditorCmd)
	}
}

func TestAsJson(t *testing.T) {
	config := NewConfig("/foo/bar", "foo", "")

	got, err := config.AsJson()

	if got == nil || err != nil {
		t.Errorf("Got '%q' Err '%q'", got, err)
	}
}

func TestSave(t *testing.T) {
	appFs := afero.NewMemMapFs()

	config := NewConfig("/foo/bar", "foo", "")

	got := config.Save(appFs)

	if got != nil {
		t.Errorf("expected 'nil' from Save. got='%q'", got)
	}
}
