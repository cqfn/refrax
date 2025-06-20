package cmd

import (
	"io"
	"os"

	"github.com/cqfn/refrax/internal/log"
	"github.com/spf13/cobra"
)

type Params struct {
	provider string
}

func Execute() error {
	return NewRootCmd(os.Stdout, os.Stderr).Execute()
}

func NewRootCmd(out io.Writer, err io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:   "refrax",
		Short: "Refrax is an AI-powered refactoring agent for Java code",
		Long:  "Refrax is an AI-powered refactoring agent for Java code. It communicates using the A2A protocol",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Set(log.NewZerolog(out, "debug"))
		},
	}
	var params Params
	root.PersistentFlags().StringVarP(&params.provider, "ai", "a", "none", "AI provider to use (e.g., openai, deepseek, none, etc.)")
	root.AddCommand(
		newRefactorCmd(&params),
		newStartCmd(&params),
	)
	return root
}
