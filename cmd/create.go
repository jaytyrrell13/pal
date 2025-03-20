package cmd

import (
	"errors"
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
	editCmd  string
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create an alias",
	Aliases: []string{"cr"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		prereqsErr := CheckCreatePrerequisites(fs)
		if prereqsErr != nil {
			return prereqsErr
		}

		cp, promptsErr := RunCreatePrompts(fs)
		if promptsErr != nil {
			return promptsErr
		}

		return RunCreateCmd(fs, cp)
	},
}

func CheckCreatePrerequisites(fs afero.Fs) error {
	ok, err := config.ConfigFileExists(fs)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Config file does not exist")
	}

	return nil
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

	cp := CreatePrompts{
		category: category,
	}

	switch category {
	case "action":
		actionRes, actionErr := ui.Input(ui.InputProps{Title: "What is the action?"})
		if actionErr != nil {
			return CreatePrompts{}, actionErr
		}

		aliasRes, aliasErr := ui.Input(ui.InputProps{Title: fmt.Sprintf("Alias for (%s).", actionRes)})
		if aliasErr != nil {
			return CreatePrompts{}, aliasErr
		}

		cp.aliases = []alias.Alias{{Name: aliasRes, Command: actionRes}}

	case "directory":
		pathRes, pathErr := ui.Input(ui.InputProps{Title: "What is the path?"})
		if pathErr != nil {
			return CreatePrompts{}, pathErr
		}

		aliasRes, aliasErr := ui.Input(ui.InputProps{Title: fmt.Sprintf("Alias for (%s).", pathRes)})
		if aliasErr != nil {
			return CreatePrompts{}, aliasErr
		}

		editCmd, editCmdErr := ui.Select("Do you want to include an edit command as well?", []huh.Option[string]{
			huh.NewOption("Yes", "yes"),
			huh.NewOption("No", "no"),
		})
		if editCmdErr != nil {
			return CreatePrompts{}, editCmdErr
		}

		cp.editCmd = editCmd
		cp.aliases = []alias.Alias{{Name: aliasRes, Command: pathRes}}

	case "parent":
		pathRes, pathErr := ui.Input(ui.InputProps{Title: "What is the path?"})
		if pathErr != nil {
			return CreatePrompts{}, pathErr
		}

		editCmd, editCmdErr := ui.Select("Do you want to include an edit command as well?", []huh.Option[string]{
			huh.NewOption("Yes", "yes"),
			huh.NewOption("No", "no"),
		})
		if editCmdErr != nil {
			return CreatePrompts{}, editCmdErr
		}

		cp.editCmd = editCmd

		files, readDirErr := afero.ReadDir(fs, pathRes)
		if readDirErr != nil {
			return CreatePrompts{}, readDirErr
		}

		var projectPaths []string
		for _, file := range files {
			if file.Name() != ".DS_Store" {
				projectPaths = append(projectPaths, pathRes+"/"+file.Name())
			}
		}

		for _, projectPath := range projectPaths {
			aliasRes, aliasErr := ui.Input(ui.InputProps{Title: fmt.Sprintf("Alias for (%s) Leave blank to skip.", projectPath)})
			if aliasErr != nil {
				return CreatePrompts{}, aliasErr
			}

			if aliasRes == "" {
				continue
			}

			cp.aliases = append(cp.aliases, alias.Alias{Name: aliasRes, Command: projectPath})
		}
	}

	return cp, nil
}

func RunCreateCmd(fs afero.Fs, cp CreatePrompts) error {
	c, configErr := config.ReadConfigFile(fs)
	if configErr != nil {
		return configErr
	}

	for _, a := range cp.aliases {
		if cp.category != "action" {
			c.Aliases = append(c.Aliases, a.ForActionCmd())
		} else {
			c.Aliases = append(c.Aliases, a)
		}

		if cp.editCmd == "yes" {
			c.Aliases = append(c.Aliases, a.ForEditCmd())
		}
	}

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
	rootCmd.AddCommand(createCmd)
}
