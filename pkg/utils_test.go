package pkg

import (
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	tempDir := t.TempDir()
	err := os.WriteFile(tempDir+"/foo", []byte("123"), 0o644)

	got := ReadFile(tempDir + "/foo")

	if got == nil || err != nil {
		t.Errorf("Expected 'nil', but got '%q'", got)
	}
}
