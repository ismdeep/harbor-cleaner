package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	"github.com/kopeisec/fp"
	"github.com/spf13/cobra"

	"github.com/ismdeep/harbor-cleaner/harbor"
)

func ReposCommand() *cobra.Command {
	var project string
	var orderBy string
	m := &cobra.Command{
		Use:   "repos",
		Short: "Show all repos",
		RunE: func(cmd *cobra.Command, args []string) error {

			type RepoInfo struct {
				Name string
				Size int64
			}

			repos, err := client.ListRepositories(project)
			if err != nil {
				return fmt.Errorf("listing repos: %w", err)
			}

			repoInfos := fp.TransformAsync(repos, func(repo string) RepoInfo {
				tags, err := client.ListRepoTags(project, repo)
				if err != nil {
					return RepoInfo{
						Name: repo,
						Size: 0,
					}
				}

				return RepoInfo{
					Name: repo,
					Size: fp.Transform(tags, func(tag harbor.TagInfo) int64 { return tag.Size }).
						Reduce(func(i1 int64, i2 int64) int64 { return i1 + i2 }, 0),
				}
			})

			if orderBy == "" {
				orderBy = "name"
			}
			switch orderBy {
			case "name":
				sort.Slice(repoInfos, func(i, j int) bool {
					return strings.Compare(repoInfos[i].Name, repoInfos[j].Name) < 0
				})
			case "size":
				sort.Slice(repoInfos, func(i, j int) bool {
					return repoInfos[i].Size < repoInfos[j].Size
				})
			default:
				return fmt.Errorf("invalid order by: %s", orderBy)
			}

			// print as table
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "IMAGE\tSIZE")
			for _, r := range repoInfos {
				_, _ = fmt.Fprintf(w, "%s\t%s\n", fmt.Sprintf("%v/%v/%v", client.Host(), project, r.Name), humanize.Bytes(uint64(r.Size)))
			}
			_ = w.Flush()

			return nil
		},
	}

	m.PersistentFlags().StringVar(&project, "project", "", "Harbor project name or ID")
	m.PersistentFlags().StringVar(&orderBy, "order", "name", "Order by, e.g. name, size")

	return m
}
