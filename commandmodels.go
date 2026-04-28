package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func CommandModels() *cobra.Command {
	var check bool
	var availableOnly bool
	var concurrency int

	cmd := &cobra.Command{
		Use:   "models",
		Short: "List available models",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if availableOnly && !check {
				return fmt.Errorf("--available-only requires --check")
			}
			if concurrency < 1 {
				return fmt.Errorf("--concurrency must be greater than 0")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := LoadFromFile(ConfigFilePath())
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			models, err := RequestModels(config.Endpoint, config.Key)
			if err != nil {
				return err
			}

			return runModelsCommand(os.Stdout, config, models, check, availableOnly, concurrency, RequestContent)
		},
	}

	cmd.Flags().BoolVar(&check, "check", false, "Check model availability with a test prompt")
	cmd.Flags().BoolVar(&availableOnly, "available-only", false, "Only show available models when used with --check")
	cmd.Flags().IntVar(&concurrency, "concurrency", 5, "Number of concurrent checks when used with --check")
	return cmd
}

func runModelsCommand(w io.Writer, config *Config, models []ModelData, check bool, availableOnly bool, concurrency int, requestContent func(endpoint string, model string, key string, input string) (string, error)) error {
	models = append([]ModelData(nil), models...)
	sort.Slice(models, func(i int, j int) bool {
		return models[i].ID < models[j].ID
	})

	statuses := map[string]string{}
	if check {
		if err := checkModelStatuses(statuses, config, models, concurrency, requestContent); err != nil {
			return err
		}
	}

	if availableOnly {
		availableModels := make([]ModelData, 0, len(models))
		for _, model := range models {
			if statuses[model.ID] == "Available" {
				availableModels = append(availableModels, model)
			}
		}
		models = availableModels
	}

	return writeModelsTable(w, models, statuses)
}

func checkModelStatuses(statuses map[string]string, config *Config, models []ModelData, concurrency int, requestContent func(endpoint string, model string, key string, input string) (string, error)) error {
	if concurrency > len(models) {
		concurrency = len(models)
	}
	if concurrency < 1 {
		concurrency = 1
	}

	jobs := make(chan ModelData)
	var wg sync.WaitGroup
	var mu sync.Mutex

	worker := func() {
		defer wg.Done()
		for model := range jobs {
			status := "Available"
			if _, err := requestContent(config.Endpoint, model.ID, config.Key, "Hi."); err != nil {
				status = "Unavailable"
			}

			mu.Lock()
			statuses[model.ID] = status
			mu.Unlock()
		}
	}

	wg.Add(concurrency)
	for range concurrency {
		go worker()
	}

	for _, model := range models {
		jobs <- model
	}
	close(jobs)
	wg.Wait()
	return nil
}

func writeModelsTable(w io.Writer, models []ModelData, statuses map[string]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if len(statuses) == 0 {
		if _, err := fmt.Fprintln(tw, "MODEL\tBY"); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintln(tw, "MODEL\tBY\tSTATUS"); err != nil {
			return err
		}
	}

	for _, model := range models {
		if len(statuses) == 0 {
			if _, err := fmt.Fprintf(tw, "%s\t%s\n", model.ID, model.OwnedBy); err != nil {
				return err
			}
			continue
		}

		status := statuses[model.ID]
		if status == "" {
			status = "Unavailable"
		}

		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\n", model.ID, model.OwnedBy, status); err != nil {
			return err
		}
	}

	return tw.Flush()
}
