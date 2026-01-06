package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	m := &cobra.Command{
		Use:           "aicc",
		Short:         "aicc",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	m.AddCommand(CommandConfig())
	m.AddCommand(CommandChat())

	if err := m.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR]: %v\n", err)
		os.Exit(1)
	}
}
