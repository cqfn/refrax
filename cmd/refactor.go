// Package cmd provides command-line interface functionality for the refrax tool.
package cmd

import (
	"strings"

	"github.com/cqfn/refrax/internal/client"
	"github.com/cqfn/refrax/internal/env"
	"github.com/cqfn/refrax/internal/log"
	"github.com/spf13/cobra"
)

func newRefactorCmd(params *Params) *cobra.Command {
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
			log.Debug("refactoring provider: %s", params.provider)
			log.Debug("project path to refactor: %s", path)
			var token string
			if params.token != "" {
				token = params.token
			} else {
				log.Info("token not provided, trying to find token in .env file")
				token = env.Token(".env", params.provider)
			}
			log.Debug("using provided token: %s...", mask(token))
			proj, err := project(path, output, params)
			if err != nil {
				return err
			}
			ref, err := client.Refactor(params.provider, token, proj, params.stats, log.Default(), params.playbook)
			log.Debug("refactor result: %s", ref)
			return err
		},
	}
	command.Flags().StringVarP(&output, "output", "o", "", "output path for the refactored code")
	return command
}

func project(path, output string, params *Params) (client.Project, error) {
	if params.mock {
		return client.NewMockProject(), nil
	}
	input := client.NewFilesystemProject(path)
	if output != "" {
		return client.NewMirrorProject(input, output)
	}
	return input, nil
}

func mask(token string) string {
	n := len(token)
	if n == 0 {
		return ""
	}
	visible := min(n, 3)
	return token[:visible] + strings.Repeat("*", n-visible)
}
