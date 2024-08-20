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
	Short:   "Something",
	Aliases: []string{"ls"},
	RunE: func(_ *cobra.Command, _ []string) error {
		return RunConfigListCmd()
	},
}

func RunConfigListCmd() error {
	AppFs := afero.NewOsFs()

	configDirPath, configDirPathErr := pkg.ConfigDirPath()
	if configDirPathErr != nil {
		return configDirPathErr
	}

	if pkg.FileMissing(AppFs, configDirPath) {
		configDirErr := pkg.MakeConfigDir(AppFs)
		if configDirErr != nil {
			return configDirErr
		}
	}

	configFilePath, configFilePathErr := pkg.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	configFile, readConfigFileErr := pkg.ReadFile(AppFs, configFilePath)
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
		{"Editor Command", c.EditorCmd},
		{"Shell", c.Shell},
		{"Extras", strings.Join(c.Extras, ", ")},
	}

	t := ui.Table(headers, rows)

	fmt.Println(t)

	return nil
}
