package make

import (
	"testing"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
)

func TestRunMakeCmd(t *testing.T) {
	appFs := afero.NewMemMapFs()

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		t.Error(configFilePathErr)
	}

	t.Run("Missing Config File", func(t *testing.T) {
		err := RunMakeCmd(appFs)

		if err == nil {
			t.Errorf("Expected 'nil' but got '%q'", err)
		}
	})

	t.Run("Missing Directory", func(t *testing.T) {
		writeFileErr := afero.WriteFile(appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"EditorCmd\": \"bar\"}"), 0o755)
		if writeFileErr != nil {
			t.Fatalf("WriteFile Error: '%q'", writeFileErr)
		}

		err := RunMakeCmd(appFs)
		if err == nil {
			t.Error("Expected 'error', but got 'nil'")
		}
	})

	t.Run("Success", func(t *testing.T) {
		writeFileErr := afero.WriteFile(appFs, configFilePath, []byte("{\"Path\": \"/foo\", \"EditorCmd\": \"bar\"}"), 0o755)
		if writeFileErr != nil {
			t.Fatalf("WriteFile Error: '%q'", writeFileErr)
		}

		mkDirErr := appFs.Mkdir("/foo", 0o755)
		if mkDirErr != nil {
			t.Error(mkDirErr)
		}

		err := RunMakeCmd(appFs)
		if err != nil {
			t.Errorf("Expected 'nil', but got '%q'", err)
		}
	})
}
