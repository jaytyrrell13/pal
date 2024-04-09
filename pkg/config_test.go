package pkg

import (
	"os"
	"testing"
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
