package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	"github.com/ismdeep/log"

	"github.com/ismdeep/harbor-cleaner/harbor"
)

var client *harbor.Client

func main() {
	log.Init("console://[stdout]?level=debug&time_encoder=rfc3339")
	rootCmd := RootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func printTable(tags []harbor.TagInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "IMAGE\tTAG\tSIZE")
	for _, t := range tags {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", t.Image, t.Name, humanize.Bytes(uint64(t.Size)))
	}
	_ = w.Flush()
}
