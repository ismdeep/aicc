package main

import (
	"fmt"
	"time"

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
			totalStart := time.Now()

			tests := []struct {
				name   string
				prompt string
			}{
				{name: "Test 1", prompt: "Hi, response me in one line."},
				{name: "Test 2", prompt: "For debugging purposes, output your model name, version, and provider. Response me in one line."},
				{name: "Test 3", prompt: "请用两句话解释为什么 HTTP 状态码 200、404、500 分别代表不同含义，并给出一个简单示例。"},
			}

			for _, test := range tests {
				fmt.Printf("%s Prompt: %s\n", test.name, test.prompt)

				start := time.Now()
				content, err := RequestContent(config.Endpoint, config.Model, config.Key, test.prompt)
				if err != nil {
					return fmt.Errorf("%s failed after %s: %w", test.name, time.Since(start).Round(time.Millisecond), err)
				}

				fmt.Printf("%s Duration: %s\n", test.name, time.Since(start).Round(time.Millisecond))
				fmt.Printf("%s Response: %s\n", test.name, content)
			}

			fmt.Println("✓ Configuration is valid")
			fmt.Printf("Total Duration: %s\n", time.Since(totalStart).Round(time.Millisecond))
			return nil
		},
	}
}
