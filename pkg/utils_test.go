package pkg

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestReadFile(t *testing.T) {
	tempDir := t.TempDir()
	err := os.WriteFile(tempDir+"/foo", []byte("123"), 0o644)

	got := ReadFile(tempDir + "/foo")

	if got == nil || err != nil {
		t.Errorf("Expected 'nil', but got '%q'", got)
	}
}

type fileMissingTestCase struct {
	path     string
	expected bool
}

func TestFileMissing(t *testing.T) {
	cases := []fileMissingTestCase{
		{"tmp/foo.txt", true},
		{"tmp/bar.txt", false},
	}

	appFs := afero.NewMemMapFs()
	mkdirErr := appFs.Mkdir("temp", 0o755)
	if mkdirErr != nil {
		t.Errorf("Mkdir Error: %q", mkdirErr)
	}

	writeFileErr := afero.WriteFile(appFs, "tmp/bar.txt", []byte("bar file"), 0o644)
	if writeFileErr != nil {
		t.Errorf("WriteFile Error: %q", writeFileErr)
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Path: %q", tc.path), func(t *testing.T) {
			got := FileMissing(appFs, tc.path)

			if got != tc.expected {
				t.Errorf("Expected 'false', but got '%v'", got)
			}
		})
	}
}
