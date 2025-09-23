package cmd

import (
	"io"
	"os"

	"github.com/cqfn/refrax/internal/client"
	"github.com/cqfn/refrax/internal/util"
	"github.com/spf13/cobra"
)

// Execute runs the root command and returns any error encountered.
func Execute() error {
	return NewRootCmd(os.Stdout, os.Stderr).Execute()
}

// NewRootCmd creates and returns the root command for Refrax.
// Command line interface for Refrax.
// @todo #2:45min Add new parameter `--tools` for constructing needed tool for the critic.
// Currently, we create Aibolit, but it would be great to have such option. Examples of such tools
// are: `aibolit`, `none`, etc. We should be able to pass multiple tools, for instance:
// `--tools=aibolit,qulice`.
func NewRootCmd(out, _ io.Writer) *cobra.Command {
	var params client.Params
	root := &cobra.Command{
		Use:              "refrax",
		Short:            "Refrax is an AI-powered refactoring agent for Java code",
		Long:             "Refrax is an AI-powered refactoring agent for Java code. It communicates using the A2A protocol",
		PersistentPreRun: func(_ *cobra.Command, _ []string) { params.Log = out },
	}
	root.PersistentFlags().StringVarP(&params.Provider, "ai", "a", "none", "AI provider to use (openai, deepseek, none)")
	root.PersistentFlags().StringVarP(&params.Token, "token", "t", "", "Token for the AI provider (if required)")
	root.PersistentFlags().StringVar(&params.Playbook, "playbook", "", "Path to a user-defined YAML playbook for AI integration")
	root.PersistentFlags().BoolVar(&params.MockProject, "mock-project", false, "Use mock project")
	root.PersistentFlags().BoolVarP(&params.Debug, "debug", "d", false, "Print debug logs")
	root.PersistentFlags().BoolVar(&params.Stats, "stats", false, "Print internal interaction statistics")
	root.PersistentFlags().StringVar(&params.Format, "stats-format", "std", "Format for statistics output (std, csv)")
	root.PersistentFlags().StringVar(&params.Soutput, "stats-output", "stats", "Output path for statistics")
	root.PersistentFlags().BoolVar(&params.Colorless, "no-colors", false, "Disable colored output")
	root.AddCommand(
		newRefactorCmd(&params),
		newStartCmd(),
	)
	root.Version = util.Version()
	root.SetVersionTemplate("refrax {{.Version}}\n")
	return root
}
