package cmd

import (
	"github.com/spf13/cobra"
)

func newStartCmd(params *Params) *cobra.Command {
	command := &cobra.Command{
		Use:     "start [agent]",
		Short:   "Starts a particular agent like fixer, critic, or facilitator",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"st"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// return facilitator.StartServer(params.provider, params.token, 8080)
			return nil
		},
	}
	return command
}
