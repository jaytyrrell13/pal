package cmd

import (
	"errors"
	"fmt"
	"maps"

	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/jaytyrrell13/pal/internal/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type UpdatedAlias struct {
	oldName  string
	newAlias alias.Alias
}

type UpdatePrompts struct {
	updatedAliases []UpdatedAlias
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update an alias",
	Aliases: []string{"up"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		prereqsErr := CheckUpdatePrerequisites(fs)
		if prereqsErr != nil {
			return prereqsErr
		}

		up, promptsErr := RunUpdatePrompts(fs)
		if promptsErr != nil {
			return promptsErr
		}

		return RunUpdateCmd(fs, up)
	},
}

func CheckUpdatePrerequisites(fs afero.Fs) error {
	configFileExists, configFileExistsErr := config.ConfigFileExists(fs)
	if configFileExistsErr != nil {
		return configFileExistsErr
	}

	if !configFileExists {
		return errors.New("Config file does not exist")
	}

	aliasesFileExists, aliasesFileExistsErr := config.AliasesFileExists(fs)
	if aliasesFileExistsErr != nil {
		return aliasesFileExistsErr
	}

	if !aliasesFileExists {
		return errors.New("Aliases file does not exist")
	}

	return nil
}

func RunUpdatePrompts(fs afero.Fs) (UpdatePrompts, error) {
	c, configErr := config.ReadConfigFile(fs)
	if configErr != nil {
		return UpdatePrompts{}, configErr
	}

	var options []huh.Option[string]
	for _, a := range c.Aliases {
		display := fmt.Sprintf("Name: '%s' Command: '%s'", a.Name, a.Command)

		options = append(options, huh.NewOption(display, a.Name))
	}

	optionsToUpdate, optionsToUpdateErr := ui.MultiSelect("Aliases to update", options)
	if optionsToUpdateErr != nil {
		return UpdatePrompts{}, optionsToUpdateErr
	}

	var updatedAliases []UpdatedAlias
	for _, a := range c.Aliases {
		for _, o := range optionsToUpdate {
			if a.Name == o {
				nameProps := ui.InputProps{
					Title: "Updated Name",
					Value: a.Name,
				}
				updatedName, updatedNameErr := ui.Input(nameProps)
				if updatedNameErr != nil {
					return UpdatePrompts{}, updatedNameErr
				}

				commandProps := ui.InputProps{
					Title: "Update Command",
					Value: a.Command,
				}
				updatedCommand, updatedCommandErr := ui.Input(commandProps)
				if updatedCommandErr != nil {
					return UpdatePrompts{}, updatedCommandErr
				}

				updatedAliases = append(updatedAliases, UpdatedAlias{
					oldName: a.Name,
					newAlias: alias.Alias{
						Name:    updatedName,
						Command: updatedCommand,
					},
				})
			}
		}
	}

	return UpdatePrompts{
		updatedAliases: updatedAliases,
	}, nil
}

func RunUpdateCmd(fs afero.Fs, up UpdatePrompts) error {
	c, configErr := config.ReadConfigFile(fs)
	if configErr != nil {
		return configErr
	}

	aliasMap := make(map[string]alias.Alias)
	for _, a := range c.Aliases {
		aliasMap[a.Name] = a
	}

	for _, ua := range up.updatedAliases {
		aliasMap[ua.oldName] = ua.newAlias
	}

	c.Aliases = []alias.Alias{}
	for v := range maps.Values(aliasMap) {
		c.Aliases = append(c.Aliases, v)
	}

	writeConfigFileErr := config.WriteConfigFile(fs, c)
	if writeConfigFileErr != nil {
		return writeConfigFileErr
	}

	return config.WriteAliasFile(fs, c)
}
