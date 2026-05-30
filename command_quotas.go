package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func QuotasCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "quotas",
		Short: "Show quota usage of all projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			quotas, err := client.Quotas()
			if err != nil {
				return fmt.Errorf("get quotas: %v", err.Error())
			}

			sort.Slice(quotas, func(i, j int) bool {
				return quotas[i].Used.Storage < quotas[j].Used.Storage
			})

			// print as table
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "PROJECT\tUSED")
			for _, quota := range quotas {
				_, _ = fmt.Fprintf(w, "%s\t%s\n", quota.Ref.Name, humanize.Bytes(uint64(quota.Used.Storage)))
			}
			_ = w.Flush()

			return nil
		},
	}

	return m
}
