package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/charmbracelet/huh"
	cfg "github.com/jaytyrrell13/pal/internal/config"
	"github.com/jaytyrrell13/pal/internal/messages"
	"github.com/jaytyrrell13/pal/internal/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type UpdatedConfig struct {
	oldKey   string
	newValue string
}

type UpdatePrompts struct {
	updatedConfig []UpdatedConfig
}

func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update pal config",
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
}

func CheckUpdatePrerequisites(fs afero.Fs) error {
	configFileExists, configFileExistsErr := cfg.ConfigFileExists(fs)
	if configFileExistsErr != nil {
		return configFileExistsErr
	}

	if !configFileExists {
		return errors.New(messages.Errors["configMissing"])
	}
	return nil
}

func RunUpdatePrompts(fs afero.Fs) (UpdatePrompts, error) {
	c, configErr := cfg.ReadConfigFile(fs)
	if configErr != nil {
		return UpdatePrompts{}, configErr
	}

	v := reflect.ValueOf(c)
	typeOfS := v.Type()

	var options []huh.Option[string]
	for i := range v.NumField() {
		fieldName := typeOfS.Field(i).Name
		fieldValue := v.Field(i).String()

		if fieldName != "Aliases" {
			display := fmt.Sprintf("Key: '%s' Value: '%s'", fieldName, fieldValue)
			options = append(options, huh.NewOption(display, fieldName))
		}
	}

	toUpdate, toUpdateErr := ui.MultiSelect("Config to update", options)
	if toUpdateErr != nil {
		return UpdatePrompts{}, toUpdateErr
	}

	var updatedConfig []UpdatedConfig
	for _, configKey := range toUpdate {
		if configKey == "Shell" {
			newShell, newShellErr := ui.Select("What shell do you use?", []huh.Option[string]{
				huh.NewOption("Bash", "bash"),
				huh.NewOption("Fish", "fish"),
				huh.NewOption("Zsh", "zsh"),
			})
			if newShellErr != nil {
				return UpdatePrompts{}, newShellErr
			}

			updatedConfig = append(updatedConfig, UpdatedConfig{
				oldKey:   configKey,
				newValue: newShell,
			})
		}
	}

	return UpdatePrompts{
		updatedConfig: updatedConfig,
	}, nil
}

func RunUpdateCmd(fs afero.Fs, up UpdatePrompts) error {
	c, configErr := cfg.ReadConfigFile(fs)
	if configErr != nil {
		return configErr
	}

	for _, v := range up.updatedConfig {
		switch v.oldKey {
		case "Shell":
			c.Shell = v.newValue
		}
	}

	return cfg.WriteConfigFile(fs, c)
}
