package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func CommandGitDiffConventionalMessage() *cobra.Command {
	return &cobra.Command{
		Use:   "gitdiffmsg",
		Short: "Git Diff Conventional Message",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := readStdin()
			if err != nil {
				return err
			}

			content = fmt.Sprintf(`
You are a tool that generates Conventional Commit messages.
Given the following git diff, summarize the change and output ONLY a commit message.

Requirements:
- Use Conventional Commit format (feat, fix, docs, refactor, perf, test, chore, ci, build, style)
- Message must be short, imperative mood.
- DO NOT include explanations.
- Output a single line.

Git diff:
%s
`, content)

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
