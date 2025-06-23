package cmd

import (
	"github.com/cqfn/refrax/internal/client"
	"github.com/cqfn/refrax/internal/log"
	"github.com/spf13/cobra"
)

func newRefactorCmd(params *Params) *cobra.Command {
	command := &cobra.Command{
		Use:     "refactor [path]",
		Short:   "refactor code in the given directory (defaults to current)",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"r"},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			log.Debug("refactoring provider: %s", params.provider)
			log.Debug("project path to refactor: %s", path)
			ref, err := client.Refactor(params.provider, params.token, project(path, params))
			log.Debug("refactor result: %s", ref)
			return err
		},
	}
	return command
}

func project(path string, params *Params) client.Project {
	if params.mock {
		return client.NewMockProject()
	}
	return client.NewFilesystemProject(path)
}
