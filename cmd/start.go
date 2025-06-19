package cmd

import (
	"github.com/cqfn/refrax/internal/facilitator"
	"github.com/spf13/cobra"
)

func newStartCmd(params *Params) *cobra.Command {
	command := &cobra.Command{
		Use:     "start",
		Aliases: []string{"st"},
		Short:   "start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return facilitator.StartServer(params.provider, 8080)
		},
	}
	return command
}
