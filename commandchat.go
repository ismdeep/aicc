package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

//go:embed prompt-git-diff-msg-en.txt
var promptGitDiffEN string

//go:embed prompt-git-diff-msg-cn.txt
var promptGitDiffCN string

func CommandGitDiffConventionalMessageEnglish() *cobra.Command {
	return &cobra.Command{
		Use:   "git-diff-msg-en",
		Short: "Git Diff Conventional Message (English)",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := readStdin()
			if err != nil {
				return err
			}

			content = fmt.Sprintf(promptGitDiffEN, content)

			config, err := LoadFromFile(ConfigFilePath())
			if err != nil {
				return err
			}

			result, err := Request(config.Endpoint, config.Model, config.Key, content)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}
}

func CommandGitDiffConventionalMessageChinese() *cobra.Command {
	return &cobra.Command{
		Use:   "git-diff-msg-cn",
		Short: "Git Diff Conventional Message (Chinese)",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := readStdin()
			if err != nil {
				return err
			}

			content = fmt.Sprintf(promptGitDiffCN, content)

			config, err := LoadFromFile(ConfigFilePath())
			if err != nil {
				return err
			}

			result, err := Request(config.Endpoint, config.Model, config.Key, content)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}
}

func CommandChat() *cobra.Command {
	return &cobra.Command{
		Use:   "chat",
		Short: "Chat",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := readStdin()
			if err != nil {
				return err
			}

			config, err := LoadFromFile(ConfigFilePath())
			if err != nil {
				return err
			}

			result, err := Request(config.Endpoint, config.Model, config.Key, content)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}
}

func readStdin() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	// Not piped?
	if info.Mode()&os.ModeCharDevice != 0 {
		return "", fmt.Errorf("no input received (stdin is empty)")
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, bufio.NewReader(os.Stdin))
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
