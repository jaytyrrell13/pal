package list

import (
	"fmt"
	"strings"

	"github.com/jaytyrrell13/pal/pkg"
	"github.com/jaytyrrell13/pal/pkg/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var ConfigListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List pal configs",
	Aliases: []string{"ls"},
	RunE: func(_ *cobra.Command, _ []string) error {
		appFs := afero.NewOsFs()

		return RunConfigListCmd(appFs)
	},
}

func RunConfigListCmd(appFs afero.Fs) error {
	configDirPath, configDirPathErr := pkg.ConfigDirPath()
	if configDirPathErr != nil {
		return configDirPathErr
	}

	if pkg.FileMissing(appFs, configDirPath) {
		configDirErr := pkg.MakeConfigDir(appFs)
		if configDirErr != nil {
			return configDirErr
		}
	}

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	configFile, readConfigFileErr := pkg.ReadFile(appFs, configFilePath)
	if readConfigFileErr != nil {
		return readConfigFileErr
	}

	c, configFileErr := pkg.FromJson(configFile)
	if configFileErr != nil {
		return configFileErr
	}

	headers := []string{"Key", "Value"}
	rows := [][]string{
		{"Path", c.Path},
		{"EditorCmd", c.EditorCmd},
		{"Shell", c.Shell},
		{"Extras", strings.Join(c.Extras, ", ")},
	}

	t := ui.Table(headers, rows)

	fmt.Println(t)

	return nil
}
