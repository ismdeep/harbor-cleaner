package main

import "github.com/spf13/cobra"

var Version string

func VersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the application",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("harbor-cleaner - Clean tags on your harbor registry")
			cmd.Println()
			cmd.Println("version:", Version)
		},
	}
}
