package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

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

	if !check {
		return writeModelsTable(w, models, nil)
	}

	return writeCheckedModelsTable(w, config, models, availableOnly, concurrency, requestContent)
}

type modelCheckResult struct {
	model  ModelData
	status string
}

func writeCheckedModelsTable(w io.Writer, config *Config, models []ModelData, availableOnly bool, concurrency int, requestContent func(endpoint string, model string, key string, input string) (string, error)) error {
	modelWidth, ownedByWidth := modelColumnWidths(models)
	if err := writeCheckedModelsHeader(w, modelWidth, ownedByWidth); err != nil {
		return err
	}

	results := make(chan modelCheckResult, len(models))
	go checkModelStatuses(results, config, models, concurrency, requestContent)

	for result := range results {
		if availableOnly && result.status != "Available" {
			continue
		}
		if err := writeCheckedModelRow(w, result.model, result.status, modelWidth, ownedByWidth); err != nil {
			return err
		}
	}
	return nil
}

func checkModelStatuses(results chan<- modelCheckResult, config *Config, models []ModelData, concurrency int, requestContent func(endpoint string, model string, key string, input string) (string, error)) {
	defer close(results)

	if concurrency > len(models) {
		concurrency = len(models)
	}
	if concurrency < 1 {
		concurrency = 1
	}

	jobs := make(chan ModelData)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for model := range jobs {
			status := "Available"
			if _, err := requestContent(config.Endpoint, model.ID, config.Key, "Hi."); err != nil {
				status = "Unavailable"
			}

			results <- modelCheckResult{model: model, status: status}
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
}

func writeModelsTable(w io.Writer, models []ModelData, statuses map[string]string) error {
	modelWidth, ownedByWidth := modelColumnWidths(models)
	if len(statuses) == 0 {
		if _, err := fmt.Fprintf(w, "%-*s  %-*s\n", modelWidth, "MODEL", ownedByWidth, "BY"); err != nil {
			return err
		}
	} else {
		if err := writeCheckedModelsHeader(w, modelWidth, ownedByWidth); err != nil {
			return err
		}
	}

	for _, model := range models {
		if len(statuses) == 0 {
			if _, err := fmt.Fprintf(w, "%-*s  %-*s\n", modelWidth, model.ID, ownedByWidth, model.OwnedBy); err != nil {
				return err
			}
			continue
		}

		status := statuses[model.ID]
		if status == "" {
			status = "Unavailable"
		}

		if err := writeCheckedModelRow(w, model, status, modelWidth, ownedByWidth); err != nil {
			return err
		}
	}

	return nil
}

func modelColumnWidths(models []ModelData) (int, int) {
	modelWidth := len("MODEL")
	ownedByWidth := len("BY")
	for _, model := range models {
		if len(model.ID) > modelWidth {
			modelWidth = len(model.ID)
		}
		if len(model.OwnedBy) > ownedByWidth {
			ownedByWidth = len(model.OwnedBy)
		}
	}
	return modelWidth, ownedByWidth
}

func writeCheckedModelsHeader(w io.Writer, modelWidth int, ownedByWidth int) error {
	_, err := fmt.Fprintf(w, "%-*s  %-*s  %s\n", modelWidth, "MODEL", ownedByWidth, "BY", "STATUS")
	return err
}

func writeCheckedModelRow(w io.Writer, model ModelData, status string, modelWidth int, ownedByWidth int) error {
	_, err := fmt.Fprintf(w, "%-*s  %-*s  %s\n", modelWidth, model.ID, ownedByWidth, model.OwnedBy, status)
	return err
}
