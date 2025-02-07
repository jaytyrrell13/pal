package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/jaytyrrell13/pal/internal/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type CreatePrompts struct {
	category string
	aliases  []alias.Alias
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		cp, err := RunCreatePrompts(fs)
		if err != nil {
			return err
		}

		return RunCreateCmd(fs, cp)
	},
}

func RunCreatePrompts(fs afero.Fs) (CreatePrompts, error) {
	category, err := ui.Select("What type of alias would you like to create?", []huh.Option[string]{
		huh.NewOption("Parent", "parent"),
		huh.NewOption("Directory", "directory"),
		huh.NewOption("Action", "action"),
	})
	if err != nil {
		return CreatePrompts{}, err
	}

	cp := CreatePrompts{}

	switch category {
	case "action":
		action, actionErr := ui.Input("What is the action?")
		if actionErr != nil {
			return CreatePrompts{}, actionErr
		}

		a, aliasErr := ui.Input(fmt.Sprintf("Alias for (%s).", action))
		if aliasErr != nil {
			return CreatePrompts{}, aliasErr
		}

		cp.aliases = []alias.Alias{{Name: a, Path: action}}

	case "directory":
		path, pathErr := ui.Input("What is the path?")
		if pathErr != nil {
			return CreatePrompts{}, pathErr
		}

		a, aliasErr := ui.Input(fmt.Sprintf("Alias for (%s).", path))
		if aliasErr != nil {
			return CreatePrompts{}, aliasErr
		}

		cp.aliases = []alias.Alias{{Name: a, Path: path}}

	case "parent":
		path, pathErr := ui.Input("What is the path?")
		if pathErr != nil {
			return CreatePrompts{}, pathErr
		}

		files, readDirErr := afero.ReadDir(fs, path)
		if readDirErr != nil {
			return CreatePrompts{}, readDirErr
		}

		var projectPaths []string
		for _, file := range files {
			if file.Name() != ".DS_Store" {
				projectPaths = append(projectPaths, path+"/"+file.Name())
			}
		}

		for _, projectPath := range projectPaths {
			a, aliasErr := ui.Input(fmt.Sprintf("Alias for (%s) Leave blank to skip.", projectPath))
			if aliasErr != nil {
				return CreatePrompts{}, aliasErr
			}

			if a == "" {
				continue
			}

			cp.aliases = append(cp.aliases, alias.Alias{Name: a, Path: projectPath})
		}
	}

	return cp, nil
}

func RunCreateCmd(fs afero.Fs, cp CreatePrompts) error {
	configFilePath, configFilePathErr := config.ConfigFilePath()
	if configFilePathErr != nil {
		return configFilePathErr
	}

	b, readFileErr := afero.ReadFile(fs, configFilePath)
	if readFileErr != nil {
		return readFileErr
	}

	var c config.Config
	jsonErr := json.Unmarshal(b, &c)
	if jsonErr != nil {
		return jsonErr
	}

	c.Aliases = append(c.Aliases, cp.aliases...)

	j, jsonErr := json.Marshal(c)
	if jsonErr != nil {
		return jsonErr
	}

	writeFileErr := afero.WriteFile(fs, configFilePath, j, 0o644)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)
}
