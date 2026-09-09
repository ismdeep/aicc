package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

func CommandVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(Version)
			return nil
		},
	}
}
