package cmd

import (
	"io"
	"os"

	"github.com/cqfn/refrax/internal/log"
	"github.com/spf13/cobra"
)

// Params holds the configuration parameters for Refrax commands.
type Params struct {
	provider string
	token    string
	playbook string
	mock     bool
	debug    bool
	stats    bool
	format   string
	soutput  string
}

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
	var params Params
	root := &cobra.Command{
		Use:   "refrax",
		Short: "Refrax is an AI-powered refactoring agent for Java code",
		Long:  "Refrax is an AI-powered refactoring agent for Java code. It communicates using the A2A protocol",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if params.debug {
				log.Set(log.NewZerolog(out, "debug"))
			} else {
				log.Set(log.NewZerolog(out, "info"))
			}
		},
	}
	root.PersistentFlags().StringVarP(&params.provider, "ai", "a", "none", "AI provider to use (e.g., openai, deepseek, none, etc.)")
	root.PersistentFlags().StringVarP(&params.token, "token", "t", "", "Token for the AI provider (if required)")
	root.PersistentFlags().StringVar(&params.playbook, "playbook", "", "Path to a user-defined YAML playbook for AI integration.")
	root.PersistentFlags().BoolVar(&params.mock, "mock", false, "Use mock project")
	root.PersistentFlags().BoolVarP(&params.debug, "debug", "d", false, "print debug logs")
	root.PersistentFlags().BoolVar(&params.stats, "stats", false, "Print internal interaction statistics")
	root.PersistentFlags().StringVar(&params.format, "stats-format", "std", "Format for statistics output (e.g., std, csv, etc.)")
	root.PersistentFlags().StringVar(&params.soutput, "stats-output", "stats", "Output path for statistics (default: stats)")
	root.AddCommand(
		newRefactorCmd(&params),
		newStartCmd(),
	)
	return root
}
