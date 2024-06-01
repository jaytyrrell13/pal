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

	writeFileErr := afero.WriteFile(appFs, "temp/foo.txt", []byte("foo file"), 0o644)
	if writeFileErr != nil {
		t.Errorf("WriteFile Error: %q", writeFileErr)
	}

	got, err := ReadFile(appFs, "temp/foo.txt")

	if got == nil || err != nil {
		t.Errorf("Expected '[]byte', but got '%q'", got)
	}
}

type appendToFileTestCase struct {
	path   string
	data   []byte
	append []byte
}

func TestAppendToFile(t *testing.T) {
	cases := []appendToFileTestCase{
		{"temp/foo.txt", []byte("one"), []byte("one")},
		{"temp/bar.txt", []byte("two"), []byte("three")},
	}

	appFs := afero.NewMemMapFs()
	mkdirErr := appFs.Mkdir("temp", 0o755)
	if mkdirErr != nil {
		t.Errorf("Mkdir Error: %q", mkdirErr)
	}

	for _, tc := range cases {
		writeFileErr := afero.WriteFile(appFs, tc.path, tc.data, 0o644)
		if writeFileErr != nil {
			t.Errorf("WriteFile Error: %q", writeFileErr)
		}

		t.Run(fmt.Sprintf("Path: %q", tc.path), func(t *testing.T) {
			got := AppendToFile(appFs, tc.path, tc.append)

			if got != nil {
				t.Errorf("Expected 'nil', but got '%v'", got)
			}
		})
	}
}

func TestRemoveFile(t *testing.T) {
	appFs := afero.NewMemMapFs()

	mkdirErr := appFs.Mkdir("temp", 0o755)
	if mkdirErr != nil {
		t.Errorf("Mkdir Error: %q", mkdirErr)
	}

	writeFileErr := afero.WriteFile(appFs, "temp/foo.txt", []byte("foo file"), 0o644)
	if writeFileErr != nil {
		t.Errorf("WriteFile Error: %q", writeFileErr)
	}

	got := RemoveFile(appFs, "temp/foo.txt")

	if got != nil {
		t.Errorf("Expected 'nil', but got '%v'", got)
	}
}

type fileMissingTestCase struct {
	path     string
	expected bool
}

func TestFileMissing(t *testing.T) {
	cases := []fileMissingTestCase{
		{"temp/foo.txt", true},
		{"temp/bar.txt", false},
	}

	appFs := afero.NewMemMapFs()
	mkdirErr := appFs.Mkdir("temp", 0o755)
	if mkdirErr != nil {
		t.Errorf("Mkdir Error: %q", mkdirErr)
	}

	writeFileErr := afero.WriteFile(appFs, "temp/bar.txt", []byte("bar file"), 0o644)
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
