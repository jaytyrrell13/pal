package set

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/spf13/afero"
)

func TestConfigSetCommand(t *testing.T) {
	t.Run("returns an error with incorrect arguments are", func(t *testing.T) {
		appFs := afero.NewMemMapFs()

		got := RunConfigSetCmd(appFs, []string{})

		if got == nil {
			t.Errorf("expected an error, but got 'nil'")
		}
	})

	pathTests := []struct {
		key   string
		value string
	}{
		{"Path", "/foobar"},
		{"path", "/foobar"},
	}

	for _, pt := range pathTests {
		t.Run(fmt.Sprintf("sets '%s' config value", pt.key), func(t *testing.T) {
			appFs := afero.NewMemMapFs()

			configFilePath, configFileErr := pkg.ConfigFilePath()
			if configFileErr != nil {
				t.Error(configFileErr)
			}

			pkg.WriteFixtureFile(t, appFs, configFilePath, configFileData())

			got := RunConfigSetCmd(appFs, []string{pt.key, pt.value})

			if got != nil {
				t.Errorf("expected 'nil', but got '%v'", got)
			}

			c, configErr := pkg.ReadConfigFile(appFs)
			if configErr != nil {
				t.Error(configErr)
			}

			if c.Path != pt.value {
				t.Errorf("expected '%s' to be '%s', but got '%s'", pt.key, pt.value, c.Path)
			}
		})
	}

	editorTests := []struct {
		key   string
		value string
	}{
		{"Editorcmd", "foobar"},
		{"editorCmd", "foobar"},
		{"editorcmd", "foobar"},
	}

	for _, et := range editorTests {
		t.Run(fmt.Sprintf("sets '%s' config value", et.key), func(t *testing.T) {
			appFs := afero.NewMemMapFs()

			configFilePath, configFileErr := pkg.ConfigFilePath()
			if configFileErr != nil {
				t.Error(configFileErr)
			}

			pkg.WriteFixtureFile(t, appFs, configFilePath, configFileData())

			got := RunConfigSetCmd(appFs, []string{et.key, et.value})

			if got != nil {
				t.Errorf("expected 'nil', but got '%v'", got)
			}

			c, configErr := pkg.ReadConfigFile(appFs)
			if configErr != nil {
				t.Error(configErr)
			}

			if c.Editorcmd != et.value {
				t.Errorf("expected '%s' to be '%s', but got '%s'", et.key, et.value, c.Editorcmd)
			}
		})
	}

	shellTests := []struct {
		key   string
		value string
	}{
		{"Shell", "Bash"},
		{"shell", "Bash"},
	}

	for _, st := range shellTests {
		t.Run(fmt.Sprintf("sets '%s' config value", st.key), func(t *testing.T) {
			appFs := afero.NewMemMapFs()

			configFilePath, configFileErr := pkg.ConfigFilePath()
			if configFileErr != nil {
				t.Error(configFileErr)
			}

			pkg.WriteFixtureFile(t, appFs, configFilePath, configFileData())

			got := RunConfigSetCmd(appFs, []string{st.key, st.value})

			if got != nil {
				t.Errorf("expected 'nil', but got '%v'", got)
			}

			c, configErr := pkg.ReadConfigFile(appFs)
			if configErr != nil {
				t.Error(configErr)
			}

			if c.Shell != st.value {
				t.Errorf("expected '%s' to be '%s', but got '%s'", st.key, st.value, c.Shell)
			}
		})
	}

	t.Run("returns an error when 'Shell' is not supported type", func(t *testing.T) {
		appFs := afero.NewMemMapFs()

		configFilePath, configFileErr := pkg.ConfigFilePath()
		if configFileErr != nil {
			t.Error(configFileErr)
		}

		pkg.WriteFixtureFile(t, appFs, configFilePath, configFileData())

		got := RunConfigSetCmd(appFs, []string{"Shell", "baz"})

		if got == nil {
			t.Error("expected an error but got 'nil'")
		}
	})

	extrasTests := []struct {
		key      string
		value    string
		expected []string
	}{
		{"Extras", "/another/one", []string{"/another/one"}},
		{"extras", "/another/one", []string{"/another/one"}},
	}

	for _, xt := range extrasTests {
		t.Run(fmt.Sprintf("sets '%s' config value", xt.key), func(t *testing.T) {
			appFs := afero.NewMemMapFs()

			configFilePath, configFileErr := pkg.ConfigFilePath()
			if configFileErr != nil {
				t.Error(configFileErr)
			}

			pkg.WriteFixtureFile(t, appFs, configFilePath, configFileData())

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

func configFileData() []byte {
	return []byte("{\"Path\": \"/foo\", \"Editorcmd\": \"bar\", \"Shell\": \"Fish\", \"Extras\": [\"/one/extra\"]}")
}
