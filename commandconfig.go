package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func ConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, ".aicc", "config.json")
}

func CommandConfigShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := LoadFromFile(ConfigFilePath())
			if err != nil {
				return err
			}

			fmt.Println(config.String())

			return nil
		},
	}
}

func CommandConfigSet() *cobra.Command {
	var endpoint string
	var model string
	var key string

	m := &cobra.Command{
		Use:   "set",
		Short: "Set configuration commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := ConfigFilePath()

			config, _ := LoadFromFile(configFile)
			if config == nil {
				config = &Config{}
			}

			if endpoint != "" {
				config.Endpoint = endpoint
			}

			if model != "" {
				config.Model = model
			}

			if key != "" {
				config.Key = key
			}

			if err := os.WriteFile(configFile, []byte(config.String()), 0644); err != nil {
				return err
			}

			return nil
		},
	}

	m.PersistentFlags().StringVar(&endpoint, "endpoint", "", "Endpoint")
	m.PersistentFlags().StringVar(&model, "model", "", "Model")
	m.PersistentFlags().StringVar(&key, "key", "", "Key")

	return m
}

func CommandConfig() *cobra.Command {
	m := &cobra.Command{
		Use:   "config",
		Short: "Configuration commands",
	}

	m.AddCommand(CommandConfigShow())
	m.AddCommand(CommandConfigSet())

	return m
}
