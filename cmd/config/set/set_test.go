package set

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
)

type testCase struct {
	key   string
	value string
}

func TestConfigSetCommand(t *testing.T) {
	t.Run("returns an error when config file is missing", func(t *testing.T) {
		appFs := afero.NewMemMapFs()

		got := RunConfigSetCmd(appFs, []string{"Path", "/foobar"})

		if got == nil {
			t.Error("expected an error but got 'nil'")
		}
	})

	argsTests := []struct {
		name string
		arg  []string
	}{
		{"no args", []string{}},
		{"one arg", []string{"one"}},
	}

	for _, at := range argsTests {
		t.Run(fmt.Sprintf("returns an error with %s", at.name), func(t *testing.T) {
			appFs := afero.NewMemMapFs()

			got := RunConfigSetCmd(appFs, at.arg)

			if got == nil {
				t.Errorf("expected an error, but got 'nil'")
			}
		})
	}

	pathTests := []testCase{
		{"Path", "/foobar"},
		{"path", "/foobar"},
		{"pATh", "/foobar"},
	}

	for _, pt := range pathTests {
		t.Run(fmt.Sprintf("sets '%s' config value", pt.key), func(t *testing.T) {
			appFs := afero.NewMemMapFs()
			makeConfigFile(t, appFs)

			got := RunConfigSetCmd(appFs, []string{pt.key, pt.value})

			if got != nil {
				t.Errorf("expected 'nil', but got '%v'", got)
			}

			assertConfigFile(t, appFs, "Path", pt)
		})
	}

	editorModeTests := []testCase{
		{"Editormode", "Skip"},
		{"editorMode", "Skip"},
		{"editormode", "Skip"},
		{"ediTOrmode", "Skip"},
		{"Editormode", "Same"},
		{"editorMode", "Same"},
		{"editormode", "Same"},
		{"ediTOrmode", "Same"},
		{"Editormode", "Unique"},
		{"editorMode", "Unique"},
		{"editormode", "Unique"},
		{"ediTOrmode", "Unique"},
	}

	for _, em := range editorModeTests {
		t.Run(fmt.Sprintf("sets '%s' config value", em.key), func(t *testing.T) {
			appFs := afero.NewMemMapFs()
			makeConfigFile(t, appFs)

			got := RunConfigSetCmd(appFs, []string{em.key, em.value})

			if got != nil {
				t.Errorf("expected 'nil', but got '%v'", got)
			}

			assertConfigFile(t, appFs, "Editormode", em)
		})
	}

	t.Run("returns an error when 'Editormode' is not supported value", func(t *testing.T) {
		appFs := afero.NewMemMapFs()
		makeConfigFile(t, appFs)

		got := RunConfigSetCmd(appFs, []string{"Editormode", "baz"})

		if got == nil {
			t.Error("expected an error but got 'nil'")
		}
	})

	editorTests := []testCase{
		{"Editorcmd", "foobar"},
		{"editorCmd", "foobar"},
		{"editorcmd", "foobar"},
		{"ediTOrcmd", "foobar"},
	}

	for _, et := range editorTests {
		t.Run(fmt.Sprintf("sets '%s' config value", et.key), func(t *testing.T) {
			appFs := afero.NewMemMapFs()
			makeConfigFile(t, appFs)

			got := RunConfigSetCmd(appFs, []string{et.key, et.value})

			if got != nil {
				t.Errorf("expected 'nil', but got '%v'", got)
			}

			assertConfigFile(t, appFs, "Editorcmd", et)
		})
	}

	shellTests := []testCase{
		{"Shell", "Bash"},
		{"shell", "Bash"},
		{"shELl", "Bash"},
	}

	for _, st := range shellTests {
		t.Run(fmt.Sprintf("sets '%s' config value", st.key), func(t *testing.T) {
			appFs := afero.NewMemMapFs()
			makeConfigFile(t, appFs)

			got := RunConfigSetCmd(appFs, []string{st.key, st.value})

			if got != nil {
				t.Errorf("expected 'nil', but got '%v'", got)
			}

			assertConfigFile(t, appFs, "Shell", st)
		})
	}

	t.Run("returns an error when 'Shell' is not supported type", func(t *testing.T) {
		appFs := afero.NewMemMapFs()
		makeConfigFile(t, appFs)

		got := RunConfigSetCmd(appFs, []string{"Shell", "baz"})

		if got == nil {
			t.Error("expected an error but got 'nil'")
		}
	})

	extrasTests := []testCase{
		{"Extras", "/another/one"},
		{"extras", "/another/one"},
		{"eXTras", "/another/one"},
	}

	for _, xt := range extrasTests {
		t.Run(fmt.Sprintf("sets '%s' config value", xt.key), func(t *testing.T) {
			appFs := afero.NewMemMapFs()
			makeConfigFile(t, appFs)

			got := RunConfigSetCmd(appFs, []string{xt.key, xt.value})

			if got != nil {
				t.Errorf("expected 'nil', but got '%v'", got)
			}

			c, configErr := pkg.ReadConfigFile(appFs)
			if configErr != nil {
				t.Error(configErr)
			}

			if !slices.Equal(c.Extras, []string{"/one/extra", "/another/one"}) {
				t.Errorf("expected 'Extras' to be equal")
			}
		})
	}
}

func assertConfigFile(t *testing.T, appFs afero.Fs, key string, test testCase) {
	t.Helper()

	c, configErr := pkg.ReadConfigFile(appFs)
	if configErr != nil {
		t.Error(configErr)
	}

	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(key)

	if f.String() != test.value {
		t.Errorf("expected '%s' to be '%s', but got '%s'", test.key, test.value, c.Shell)
	}
}

func makeConfigFile(t *testing.T, appFs afero.Fs) {
	t.Helper()

	configFilePath, configFileErr := pkg.ConfigFilePath()
	if configFileErr != nil {
		t.Error(configFileErr)
	}

	data := []byte("{\"Path\": \"/foo\", \"Editormode\": \"Same\", \"Editorcmd\": \"bar\", \"Shell\": \"Fish\", \"Extras\": [\"/one/extra\"]}")

	pkg.WriteFixtureFile(t, appFs, configFilePath, data)
}
