package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CommandTest() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test configuration by sending a test message",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := LoadFromFile(ConfigFilePath())
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Println("Testing configuration...")
			content, err := RequestContent(config.Endpoint, config.Model, config.Key, "hello")
			if err != nil {
				return fmt.Errorf("test failed: %w", err)
			}

			fmt.Println("✓ Configuration is valid")
			fmt.Printf("Response: %s\n", content)
			return nil
		},
	}
}
