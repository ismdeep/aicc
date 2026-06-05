package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	m := &cobra.Command{
		Use:   "aicc",
		Short: "aicc",
		Long: `aicc

Environment variables:
  OPENAI_BASE_URL    Override the API endpoint (config: endpoint)
  OPENAI_API_KEY     Override the API key (config: key)
  OPENAI_MODEL       Override the model name (config: model)

Environment variables take precedence over values in ~/.aicc/config.json.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	m.AddCommand(CommandConfig())
	m.AddCommand(CommandChat())
	m.AddCommand(CommandModels())
	m.AddCommand(CommandTest())

	if err := m.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR]: %v\n", err)
		os.Exit(1)
	}
}
