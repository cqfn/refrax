package cmd

import (
	"io"
	"os"

	"github.com/cqfn/refrax/internal/log"
	"github.com/spf13/cobra"
)

type Params struct {
	provider string
	token    string
	mock     bool
	debug    bool
	stats    bool
}

func Execute() error {
	return NewRootCmd(os.Stdout, os.Stderr).Execute()
}

func NewRootCmd(out io.Writer, err io.Writer) *cobra.Command {
	var params Params
	root := &cobra.Command{
		Use:   "refrax",
		Short: "Refrax is an AI-powered refactoring agent for Java code",
		Long:  "Refrax is an AI-powered refactoring agent for Java code. It communicates using the A2A protocol",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if params.debug {
				log.Set(log.NewZerolog(out, "debug"))
			}
		},
	}
	root.PersistentFlags().StringVarP(&params.provider, "ai", "a", "none", "AI provider to use (e.g., openai, deepseek, none, etc.)")
	root.PersistentFlags().StringVarP(&params.token, "token", "t", "", "Token for the AI provider (if required)")
	root.PersistentFlags().BoolVar(&params.mock, "mock", false, "Use mock project")
	root.PersistentFlags().BoolVarP(&params.debug, "debug", "d", false, "print debug logs")
	root.PersistentFlags().BoolVar(&params.stats, "stats", false, "Print internal interaction statistics")
	root.AddCommand(
		newRefactorCmd(&params),
		newStartCmd(&params),
	)
	return root
}
