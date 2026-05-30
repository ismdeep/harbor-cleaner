package main

import (
	"context"
	"errors"

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

			if cmd.Use == "version" {
				return nil
			}

			// check
			if endpoint == "" {
				return errors.New("endpoint is required")
			}
			if username == "" {
				return errors.New("username is required")
			}
			if password == "" {
				return errors.New("password is required")
			}
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

	m.AddCommand(CleanCommand())
	m.AddCommand(TagsCommand())
	m.AddCommand(VersionCommand())
	m.AddCommand(ReposCommand())

	return m
}
