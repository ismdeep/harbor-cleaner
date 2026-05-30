package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func TagsCommand() *cobra.Command {
	var project string
	m := &cobra.Command{
		Use:   "tags",
		Short: "Show all tags",
		RunE: func(cmd *cobra.Command, args []string) error {

			allTags, err := client.ListProjectTags(project)
			if err != nil {
				return fmt.Errorf("listing tags: %w", err)
			}

			printTable(allTags)

			return nil
		},
	}

	m.PersistentFlags().StringVar(&project, "project", "", "Harbor project name or ID")

	return m
}
