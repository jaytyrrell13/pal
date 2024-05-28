package pkg

import (
	"fmt"
	"testing"
)

func TestAliasFilePath(t *testing.T) {
	configDir, configDirErr := ConfigDirPath()
	if configDirErr != nil {
		t.Error(configDirErr)
	}

	got, err := AliasFilePath()
	expected := configDir + "/aliases"

	if got != expected || err != nil {
		t.Errorf("Expected '%q', but got '%q'", expected, got)
	}
}

type testCase struct {
	config   Config
	expected string
}

func TestMakeAliasCommands(t *testing.T) {
	cases := []testCase{
		{Config{}, "alias foo=\"cd /foo/bar\"\n"},
		{Config{EditorCmd: "nvim"}, "alias foo=\"cd /foo/bar\"\nalias efoo=\"cd /foo/bar && nvim\"\n"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("EditorCmd: %q", tc.config.EditorCmd), func(t *testing.T) {
			got := MakeAliasCommands("foo", "/foo/bar", tc.config)

			if got != tc.expected {
				t.Errorf("Expected '%q', but got '%q'", tc.expected, got)
			}
		})
	}
}
