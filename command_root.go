package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ismdeep/harbor-cleaner/harbor"
)

func RootCommand() *cobra.Command {
	var endpoint, username, password string
	var concurrency int

	m := &cobra.Command{
		Use:   "harbor-cleaner",
		Short: "Delete Harbor image tags matching regex filters",
		Long:  "Connect to a Harbor registry, list image tags matching the given regex filters, and delete them after confirmation.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			if concurrency < 1 {
				concurrency = 1
			}

			var err error
			client, err = harbor.NewClient(ctx, endpoint, username, password, concurrency)
			if err != nil {
				return err
			}

			return nil
		},
	}

	m.PersistentFlags().StringVar(&endpoint, "endpoint", "", "Harbor API endpoint (e.g. https://docker.example.com)")
	m.PersistentFlags().StringVar(&username, "username", "", "Harbor username")
	m.PersistentFlags().StringVar(&password, "password", "", "Harbor password")
	m.PersistentFlags().IntVar(&concurrency, "concurrency", 4, "Number of concurrent delete operations")
	_ = m.MarkPersistentFlagRequired("endpoint")
	_ = m.MarkPersistentFlagRequired("username")
	_ = m.MarkPersistentFlagRequired("password")

	m.AddCommand(CleanCommand())
	m.AddCommand(TagsCommand())

	return m
}
