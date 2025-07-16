// Package cmd provides command-line interface functionality for the refrax tool.
package cmd

import (
	"github.com/cqfn/refrax/internal/client"
	"github.com/spf13/cobra"
)

func newRefactorCmd(params *client.Params) *cobra.Command {
	var output string
	command := &cobra.Command{
		Use:     "refactor [path]",
		Short:   "refactor code in the given directory (defaults to current)",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"r"},
		RunE: func(_ *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			params.Input = path
			params.Output = output
			_, err := client.Refactor(params)
			return err
		},
	}
	command.Flags().StringVarP(&output, "output", "o", "", "output path for the refactored code")
	return command
}
