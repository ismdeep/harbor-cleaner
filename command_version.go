package main

import (
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var Version string

func VersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the application",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("harbor-cleaner - Clean tags on your harbor registry")
			cmd.Println()
			cmd.Printf("version:      %s\n", Version)
			cmd.Printf("go:           %s\n", runtime.Version())
			cmd.Printf("os:           %s\n", runtime.GOOS)
			cmd.Printf("arch:         %s\n", runtime.GOARCH)
			if bi, ok := debug.ReadBuildInfo(); ok {
				for _, s := range bi.Settings {
					switch s.Key {
					case "vcs.revision":
						cmd.Printf("vcs.revision: %s\n", s.Value)
					case "vcs.time":
						cmd.Printf("vcs.time:     %s\n", s.Value)
					}
				}
			}
		},
	}
}
