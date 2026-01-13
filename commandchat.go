package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func CommandChat() *cobra.Command {
	var prompt string

	m := &cobra.Command{
		Use:   "chat",
		Short: "Chat",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := LoadFromFile(ConfigFilePath())
			if err != nil {
				return err
			}

			content, err := prepareContent(config.GetPromptDir(), prompt)
			if err != nil {
				return err
			}

			result, err := RequestContent(config.Endpoint, config.Model, config.Key, content)
			if err != nil {
				return err
			}

			writeShell := filepath.Join(config.GetPromptDir(), prompt, "write.sh")

			stat, err := os.Stat(writeShell)
			if err == nil && stat.Mode().IsRegular() {
				writeCmd := exec.Command("bash", writeShell)
				writeCmd.Stdin = bytes.NewBufferString(result)
				writeCmd.Stdout = os.Stdout
				writeCmd.Stderr = os.Stderr
				return writeCmd.Run()
			}

			fmt.Println(result)

			return nil
		},
	}

	m.PersistentFlags().StringVarP(&prompt, "prompt", "p", "", "prompt")

	return m
}

func prepareContent(promptDir string, prompt string) (string, error) {
	var promptContent string
	readFromStdin := true

	if prompt != "" {
		promptContentRaw, err := os.ReadFile(
			filepath.Join(
				promptDir, prompt, "prompt.txt"))
		if err != nil {
			return "", fmt.Errorf("failed to read prompt file: %w", err)
		}
		promptContent = string(promptContentRaw)

		prePromptShellPath := filepath.Join(
			promptDir, prompt, "read.sh")
		stat, err := os.Stat(prePromptShellPath)
		if err == nil && stat.Mode().IsRegular() {
			preCmd := exec.Command("bash", prePromptShellPath)
			output, _ := preCmd.CombinedOutput()
			promptContent += string(output)
			readFromStdin = false
		}
	}

	if readFromStdin {
		stdinContent, err := readStdin()
		if err != nil {
			return "", err
		}
		return promptContent + stdinContent, nil
	}

	return promptContent, nil
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
