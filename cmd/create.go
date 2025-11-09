package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/jaytyrrell13/pal/internal/alias"
	"github.com/jaytyrrell13/pal/internal/config"
	"github.com/jaytyrrell13/pal/internal/messages"
	"github.com/jaytyrrell13/pal/internal/ui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type CreatePrompts struct {
	category string
	aliases  []alias.Alias
	editCmd  string
}

func NewCreateCmd() *cobra.Command {
	return &cobra.Command{
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
}

func CheckCreatePrerequisites(fs afero.Fs) error {
	ok, err := config.ConfigFileExists(fs)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New(messages.Errors["configMissing"])
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
		cp, err = createActionAlias(cp)
		if err != nil {
			return CreatePrompts{}, err
		}
	case "directory":
		cp, err = createDirectoryAlias(cp)
		if err != nil {
			return CreatePrompts{}, err
		}

	case "parent":
		cp, err = createParentAlias(fs, cp)
		if err != nil {
			return CreatePrompts{}, err
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
			c.Aliases = append(c.Aliases, a.ToActionCmd())
		} else {
			c.Aliases = append(c.Aliases, a)
		}

		if cp.editCmd == "yes" {
			c.Aliases = append(c.Aliases, a.ToEditCmd())
		}
	}

	writeConfigFileErr := config.WriteConfigFile(fs, c)
	if writeConfigFileErr != nil {
		return writeConfigFileErr
	}

	return config.WriteAliasFile(fs, c)
}

func createActionAlias(cp CreatePrompts) (CreatePrompts, error) {
	actionRes, actionErr := ui.Input(ui.InputProps{Title: "What is the action?"})
	if actionErr != nil {
		return CreatePrompts{}, actionErr
	}

	aliasRes, aliasErr := ui.Input(ui.InputProps{Title: fmt.Sprintf("Alias for (%s).", actionRes)})
	if aliasErr != nil {
		return CreatePrompts{}, aliasErr
	}

	cp.aliases = []alias.Alias{{Name: aliasRes, Command: actionRes}}

	return cp, nil
}

func createDirectoryAlias(cp CreatePrompts) (CreatePrompts, error) {
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

	expandedPath, expandedPathErr := expandPath(pathRes)
	if expandedPathErr != nil {
		return CreatePrompts{}, expandedPathErr
	}

	cp.aliases = []alias.Alias{{Name: aliasRes, Command: expandedPath}}

	return cp, nil
}

func createParentAlias(fs afero.Fs, cp CreatePrompts) (CreatePrompts, error) {
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

	expandedPath, expandedPathErr := expandPath(pathRes)
	if expandedPathErr != nil {
		return CreatePrompts{}, expandedPathErr
	}

	cp.editCmd = editCmd

	files, readDirErr := afero.ReadDir(fs, expandedPath)
	if readDirErr != nil {
		return CreatePrompts{}, readDirErr
	}

	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}

		projectPath := expandedPath + "/" + file.Name()

		aliasRes, aliasErr := ui.Input(ui.InputProps{Title: fmt.Sprintf("Alias for (%s) Leave blank to skip.", projectPath)})
		if aliasErr != nil {
			return CreatePrompts{}, aliasErr
		}

		if aliasRes == "" {
			continue
		}

		cp.aliases = append(cp.aliases, alias.Alias{Name: aliasRes, Command: projectPath})
	}

	return cp, nil
}

func expandPath(path string) (string, error) {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", homeErr
	}

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	}

	return path, nil
}
