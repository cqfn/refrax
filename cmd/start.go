package cmd

import (
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "start [agent]",
		Short:   "Start a particular agent (fixer, critic, facilitator)",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"st"},
		RunE: func(_ *cobra.Command, _ []string) error {
			panic("Start command is not implemented yet")
		},
	}
	return command
}
