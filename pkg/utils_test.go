package pkg

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
)

func TestReadFile(t *testing.T) {
	appFs := afero.NewMemMapFs()

	mkdirErr := appFs.Mkdir("temp", 0o755)
	if mkdirErr != nil {
		t.Errorf("Mkdir Error: %q", mkdirErr)
	}

	writeFileErr := afero.WriteFile(appFs, "tmp/foo.txt", []byte("foo file"), 0o644)
	if writeFileErr != nil {
		t.Errorf("WriteFile Error: %q", writeFileErr)
	}

	got, err := ReadFile(appFs, "tmp/foo.txt")

	if got == nil || err != nil {
		t.Errorf("Expected '[]byte', but got '%q'", got)
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
				t.Errorf("Expected '%v', but got '%v'", tc.expected, got)
			}
		})
	}
}
