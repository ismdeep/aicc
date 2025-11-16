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

	return filepath.Join(homeDir, ".aicc.json")
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

			config := Config{
				Endpoint: endpoint,
				Model:    model,
				Key:      key,
			}

			configFile := ConfigFilePath()

			if err := os.WriteFile(configFile, []byte(config.String()), 0644); err != nil {
				return err
			}

			return nil
		},
	}

	m.PersistentFlags().StringVar(&endpoint, "endpoint", "", "Endpoint")
	m.PersistentFlags().StringVar(&model, "model", "", "Model")
	m.PersistentFlags().StringVar(&key, "key", "", "Key")

	_ = m.MarkPersistentFlagRequired("endpoint")
	_ = m.MarkPersistentFlagRequired("model")
	_ = m.MarkPersistentFlagRequired("key")

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
