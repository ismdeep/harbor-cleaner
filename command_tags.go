package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func TagsCommand() *cobra.Command {
	var project string
	var orderBy string
	m := &cobra.Command{
		Use:   "tags",
		Short: "Show all tags",
		RunE: func(cmd *cobra.Command, args []string) error {

			allTags, err := client.ListProjectTags(project)
			if err != nil {
				return fmt.Errorf("listing tags: %w", err)
			}

			if orderBy == "" {
				orderBy = "name"
			}
			switch orderBy {
			case "name":
				sort.Slice(allTags, func(i, j int) bool {
					return strings.Compare(allTags[i].Image+allTags[i].Name, allTags[j].Image+allTags[j].Name) < 0
				})
			case "size":
				sort.Slice(allTags, func(i, j int) bool {
					return allTags[i].Size < allTags[j].Size
				})
			default:
				return fmt.Errorf("invalid order by: %s", orderBy)
			}

			printTable(allTags)

			return nil
		},
	}

	m.PersistentFlags().StringVar(&project, "project", "", "Harbor project name or ID")
	m.PersistentFlags().StringVar(&orderBy, "order", "name", "Order by, e.g. name, size")

	return m
}
