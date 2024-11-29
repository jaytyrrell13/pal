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

func TestConfigFileMissing(t *testing.T) {
	t.Run("when config file is present", func(t *testing.T) {
		appFs := afero.NewMemMapFs()

		configFilePath, configFileErr := ConfigFilePath()
		if configFileErr != nil {
			t.Error(configFileErr)
		}

		WriteFixtureFile(t, appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"Editorcmd\": \"bar\"}"))

		got, err := ConfigFileMissing(appFs)
		if err != nil {
			t.Errorf("expected 'nil' but got '%q'", err)
		}

		if got != false {
			t.Errorf("expected 'false' but got '%v'", got)
		}
	})

	t.Run("when config file is missing", func(t *testing.T) {
		appFs := afero.NewMemMapFs()

		got, err := ConfigFileMissing(appFs)
		if err != nil {
			t.Errorf("expected 'nil' but got '%q'", err)
		}

		if got != true {
			t.Errorf("expected 'true' but got '%v'", got)
		}
	})
}

func TestReadConfigFile(t *testing.T) {
	t.Run("when config file is present", func(t *testing.T) {
		appFs := afero.NewMemMapFs()

		configFilePath, configFileErr := ConfigFilePath()
		if configFileErr != nil {
			t.Error(configFileErr)
		}

		WriteFixtureFile(t, appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"Editorcmd\": \"bar\"}"))

		got, err := ReadConfigFile(appFs)
		if err != nil {
			t.Errorf("expected 'nil' but got '%q'", err)
		}

		if got.Path != "/foo" {
			t.Errorf("expected Path to be '/foo' but got '%s'", got.Path)
		}

		if got.Editorcmd != "bar" {
			t.Errorf("expected Editorcmd to be 'bar' but got '%s'", got.Editorcmd)
		}
	})

	t.Run("when config file is missing", func(t *testing.T) {
		appFs := afero.NewMemMapFs()

		got, err := ReadConfigFile(appFs)

		if err == nil {
			t.Errorf("expected an error but got 'nil'")
		}

		if got.Path != "" {
			t.Errorf("expected Path to be empty but got '%s'", got.Path)
		}

		if got.Editorcmd != "" {
			t.Errorf("expected Editorcmd to be empty but got '%s'", got.Path)
		}
	})
}

func TestSaveExtraDir(t *testing.T) {
	appFs := afero.NewMemMapFs()

	configFilePath, configFileErr := ConfigFilePath()
	if configFileErr != nil {
		t.Error(configFileErr)
	}

	WriteFixtureFile(t, appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"Editorcmd\": \"bar\"}"))

	got := SaveExtraDir(appFs, "/bar/baz")

	if got != nil {
		t.Errorf("Expected 'nil', but got '%q'", got)
	}
}

func TestFromJson(t *testing.T) {
	json := "{\"Path\": \"/foo/bar\", \"Editorcmd\": \"foo\", \"Extras\":null}"
	got, err := FromJson([]byte(json))

	expected, configErr := NewConfig("/foo/bar", SameEditorMode, "foo", "")

	if configErr != nil {
		t.Errorf("Expected 'nil', but got '%q'", configErr)
	}

	if got.Path != expected.Path || got.Editorcmd != expected.Editorcmd || err != nil {
		t.Errorf("Expected Path '%q' Editorcmd '%q', but got Path '%q' Editorcmd '%q'", expected.Editorcmd, expected.Path, got.Path, got.Editorcmd)
	}
}

func TestAsJson(t *testing.T) {
	config, configErr := NewConfig("/foo/bar", SameEditorMode, "foo", "")

	if configErr != nil {
		t.Errorf("Expected 'nil', but got '%q'", configErr)
	}

	got, err := config.AsJson()

	if got == nil || err != nil {
		t.Errorf("Got '%q' Err '%q'", got, err)
	}
}

func TestSave(t *testing.T) {
	appFs := afero.NewMemMapFs()

	config, configErr := NewConfig("/foo/bar", SameEditorMode, "foo", "")

	if configErr != nil {
		t.Errorf("Expected 'nil', but got '%q'", configErr)
	}

	got := config.Save(appFs)

	if got != nil {
		t.Errorf("expected 'nil' from Save. got='%q'", got)
	}
}
