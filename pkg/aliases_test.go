package pkg

import (
	"os"
	"testing"
)

func TestAliasFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()

	got := AliasFilePath()
	expected := homeDir + "/.pal"

	if got != expected || err != nil {
		t.Fatalf("TestAliasFilePath: %q", got)
	}
}

func TestMakeAliasCommandsWithoutEditorCmd(t *testing.T) {
	output := MakeAliasCommands("foo", "/foo/bar", Config{})

	if output != "alias foo=\"cd /foo/bar\"\n" {
		t.Fatalf("TestMakeAliasCommands: %q", output)
	}
}

func TestMakeAliasCommandsWithEditorCmd(t *testing.T) {
	config := Config{
		EditorCmd: "nvim",
	}
	output := MakeAliasCommands("foo", "/foo/bar", config)

	expected := "alias foo=\"cd /foo/bar\"\nalias efoo=\"cd /foo/bar && nvim\"\n"
	if output != expected {
		t.Fatalf("TestMakeAliasCommands: %q Expected: %q", output, expected)
	}
}
