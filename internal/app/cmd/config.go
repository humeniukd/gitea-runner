// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gitea.com/gitea/runner/internal/pkg/config"

	"github.com/spf13/cobra"
)

// loadConfigCmd edits an existing config file in place, keeping its comments
// and layout, so settings can be changed from provisioning scripts.
func loadConfigCmd(configFile *string) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Generate, read and edit config files",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	configCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Print an example config file",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s", config.Example)
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Print the value of a config key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := resolveConfigFile(cmd, configFile)
			if err != nil {
				return err
			}
			value, err := config.GetValue(file, args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), value)
			return nil
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>...",
		Short: "Set the value of a config key",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := resolveConfigFile(cmd, configFile)
			if err != nil {
				return err
			}
			return config.SetValue(file, args[0], args[1:]...)
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "add <key> <value>...",
		Short: "Append values to a list config key",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := resolveConfigFile(cmd, configFile)
			if err != nil {
				return err
			}
			return config.AddValue(file, args[0], args[1:]...)
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "remove <key> <value>...",
		Short: "Remove values from a list config key",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := resolveConfigFile(cmd, configFile)
			if err != nil {
				return err
			}
			return config.RemoveValue(file, args[0], args[1:]...)
		},
	})

	return configCmd
}

// defaultConfigFileNames are looked up in the working directory and next to the
// executable when --config is not given.
var defaultConfigFileNames = []string{"config.yaml", "config.yml"}

func resolveConfigFile(cmd *cobra.Command, configFile *string) (string, error) {
	if configFile != nil && *configFile != "" {
		return *configFile, nil
	}

	var dirs []string
	if wd, err := os.Getwd(); err == nil {
		dirs = append(dirs, wd)
	}
	if exe, err := os.Executable(); err == nil {
		if dir := filepath.Dir(exe); !slices.Contains(dirs, dir) {
			dirs = append(dirs, dir)
		}
	}

	for _, dir := range dirs {
		for _, name := range defaultConfigFileNames {
			candidate := filepath.Join(dir, name)
			if stat, err := os.Stat(candidate); err == nil && !stat.IsDir() {
				fmt.Fprintf(cmd.ErrOrStderr(), "using config file %q\n", candidate)
				return candidate, nil
			}
		}
	}

	return "", fmt.Errorf("no %s found in %s, pass one with --config",
		strings.Join(defaultConfigFileNames, " or "), strings.Join(dirs, " or "))
}
