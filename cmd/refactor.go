package cmd

import (
	"fmt"

	"github.com/cqfn/refrax/internal/client"
	"github.com/cqfn/refrax/internal/log"
	"github.com/spf13/cobra"
)

func newRefactorCmd(params *Params) *cobra.Command {
	command := &cobra.Command{
		Use:     "refactor",
		Aliases: []string{"r"},
		Short:   "refactor code",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("refactor command called with provider: %s", params.provider)
			ref, err := client.Refactor("code refactoring...")
			log.Debug("refactor result: %s", ref)
			fmt.Println("Refrax says:", ref)
			return err
		},
	}
	return command
}
