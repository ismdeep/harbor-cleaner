package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/ismdeep/log"
	"github.com/kopeisec/fp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ismdeep/harbor-cleaner/harbor"
)

func main() {
	var endpoint, username, password, project string
	var filters []string
	var concurrency int

	rootCmd := &cobra.Command{
		Use:   "harbor-cleaner",
		Short: "Delete Harbor image tags matching regex filters",
		Long:  "Connect to a Harbor registry, list image tags matching the given regex filters, and delete them after confirmation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			if len(filters) == 0 {
				return fmt.Errorf("at least one --filter regex is required")
			}
			var compiledFilters []*regexp.Regexp
			for _, f := range filters {
				re, err := regexp.Compile(f)
				if err != nil {
					return fmt.Errorf("invalid regex %q: %w", f, err)
				}
				compiledFilters = append(compiledFilters, re)
			}
			filterFunc := func(filterRegexpList []*regexp.Regexp) func(tag harbor.TagInfo) bool {
				return func(tag harbor.TagInfo) bool {
					for _, re := range filterRegexpList {
						if re.MatchString(tag.Name) {
							return true
						}
					}
					return false
				}
			}

			client, err := harbor.NewClient(ctx, endpoint, username, password)
			if err != nil {
				return err
			}

			allTags, err := client.ListProjectTags(project)
			if err != nil {
				return fmt.Errorf("listing tags: %w", err)
			}

			deleteTags := fp.Wrap(allTags).Filter(filterFunc(compiledFilters))
			if len(deleteTags) == 0 {
				log.WithContext(ctx).Info("No matching tags found. Nothing to delete.")
				return nil
			}

			printTable(deleteTags)

			fmt.Printf("\nThere are %d tag(s), Found %d tag(s) matching filter(s). Delete? [y/N] ", len(allTags), len(deleteTags))
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" {
				fmt.Println("Cancelled.")
				return nil
			}

			var wg sync.WaitGroup
			sem := make(chan struct{}, concurrency)
			total := len(deleteTags)
			for i, tag := range deleteTags {
				wg.Add(1)
				sem <- struct{}{}
				go func(idx int, t harbor.TagInfo) {
					defer wg.Done()
					defer func() { <-sem }()
					log.WithContext(ctx).Info("Deleting tag", zap.String("progress", fmt.Sprintf("%v/%v", i+1, total)), zap.String("image", t.Image), zap.String("tag", t.Name))
					if err := client.DeleteTag(project, t); err != nil {
						log.WithContext(ctx).Error("Deleting tag", zap.String("progress", fmt.Sprintf("%v/%v", i+1, total)), zap.String("image", t.Image), zap.String("tag", t.Name), zap.Error(err))
					}
				}(i, tag)
			}
			wg.Wait()

			log.WithContext(ctx).Info("Finished")
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "Harbor API endpoint (e.g. https://docker.example.com)")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Harbor username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Harbor password")
	rootCmd.PersistentFlags().StringVar(&project, "project", "", "Harbor project name or ID")
	rootCmd.PersistentFlags().StringArrayVarP(&filters, "filter", "f", nil, "Regex filter for tag names (can be specified multiple times)")
	rootCmd.PersistentFlags().IntVar(&concurrency, "concurrency", 4, "Number of concurrent delete operations")
	_ = rootCmd.MarkPersistentFlagRequired("endpoint")
	_ = rootCmd.MarkPersistentFlagRequired("project")
	_ = rootCmd.MarkPersistentFlagRequired("username")
	_ = rootCmd.MarkPersistentFlagRequired("password")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func printTable(tags []harbor.TagInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "IMAGE\tTAG")
	for _, t := range tags {
		_, _ = fmt.Fprintf(w, "%s\t%s\n", t.Image, t.Name)
	}
	_ = w.Flush()
}
