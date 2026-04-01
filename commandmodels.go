package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func CommandModels() *cobra.Command {
	return &cobra.Command{
		Use:   "models",
		Short: "List available models",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := LoadFromFile(ConfigFilePath())
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			models, err := RequestModels(config.Endpoint, config.Key)
			if err != nil {
				return err
			}

			return writeModelsTable(os.Stdout, models)
		},
	}
}

func writeModelsTable(w io.Writer, models []ModelData) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "MODEL\tBY"); err != nil {
		return err
	}

	for _, model := range models {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", model.ID, model.OwnedBy); err != nil {
			return err
		}
	}

	return tw.Flush()
}
