package cmd

import (
	"errors"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/jaytyrrell13/pal/internal/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type RemovePrompts struct {
	aliasesToRemove []string
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an alias",
	Aliases: []string{"rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		prereqsErr := CheckRemovePrerequisites(fs)
		if prereqsErr != nil {
			return prereqsErr
		}

		rp, promptsErr := RunRemovePrompts(fs)
		if promptsErr != nil {
			return promptsErr
		}

		return RunRemoveCmd(fs, rp)
	},
}

func CheckRemovePrerequisites(fs afero.Fs) error {
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

func RunRemovePrompts(fs afero.Fs) (RemovePrompts, error) {
	c, configErr := config.ReadConfigFile(fs)
	if configErr != nil {
		return RemovePrompts{}, configErr
	}

	var options []huh.Option[string]
	for _, a := range c.Aliases {
		options = append(options, huh.NewOption(a.Name, a.Name))
	}

	optionsToRemove, optionsToRemoveErr := ui.MultiSelect("Aliases to remove", options)
	if optionsToRemoveErr != nil {
		return RemovePrompts{}, optionsToRemoveErr
	}

	return RemovePrompts{
		aliasesToRemove: optionsToRemove,
	}, nil
}

func RunRemoveCmd(fs afero.Fs, dp RemovePrompts) error {
	c, configErr := config.ReadConfigFile(fs)
	if configErr != nil {
		return configErr
	}

	var aliasesToKeep []alias.Alias
	for _, a := range c.Aliases {
		if !slices.Contains(dp.aliasesToRemove, a.Name) {
			aliasesToKeep = append(aliasesToKeep, a)
		}
	}

	c.Aliases = aliasesToKeep

	writeConfigFileErr := config.WriteConfigFile(fs, c)
	if writeConfigFileErr != nil {
		return writeConfigFileErr
	}

	writeAliasFileErr := config.WriteAliasFile(fs, c)
	if writeAliasFileErr != nil {
		return writeAliasFileErr
	}

	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
