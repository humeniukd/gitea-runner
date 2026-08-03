// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"errors"
	"fmt"

	"gitea.com/gitea/runner/internal/pkg/config"

	"github.com/spf13/cobra"
)

// loadConfigCmd edits an existing config file in place, keeping its comments
// and layout, so settings can be changed from provisioning scripts.
func loadConfigCmd(configFile *string) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Read and edit a config file",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	configCmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Print the value of a config key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireConfigFile(configFile); err != nil {
				return err
			}
			value, err := config.GetValue(*configFile, args[0])
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
		RunE: func(_ *cobra.Command, args []string) error {
			if err := requireConfigFile(configFile); err != nil {
				return err
			}
			return config.SetValue(*configFile, args[0], args[1:]...)
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "add <key> <value>...",
		Short: "Append values to a list config key",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := requireConfigFile(configFile); err != nil {
				return err
			}
			return config.AddValue(*configFile, args[0], args[1:]...)
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "remove <key> <value>...",
		Short: "Remove values from a list config key",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := requireConfigFile(configFile); err != nil {
				return err
			}
			return config.RemoveValue(*configFile, args[0], args[1:]...)
		},
	})

	return configCmd
}

func requireConfigFile(configFile *string) error {
	if configFile == nil || *configFile == "" {
		return errors.New("no config file given, pass one with --config")
	}
	return nil
}
