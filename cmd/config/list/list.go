package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jaytyrrell13/pal/pkg"
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

	rows := [][]string{
		{"Path", c.Path},
		{"Editor Command", c.EditorCmd},
		{"Shell", c.Shell},
		{"Extras", strings.Join(c.Extras, ", ")},
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("245"))).
		Headers("Key", "Value").
		StyleFunc(func(row, col int) lipgloss.Style {
			baseStyle := lipgloss.NewStyle().Padding(0, 2)

			switch {
			case row == 0:
				return baseStyle.Bold(true)
			case row%2 == 0:
				return baseStyle.Foreground(lipgloss.Color("240"))
			default:
				return baseStyle
			}
		}).
		Rows(rows...)

	fmt.Println(t)

	return nil
}
