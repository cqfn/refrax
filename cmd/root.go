package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "refrax",
		Short: "Refrax is an AI-powered refactoring agent for Java code",
		Long:  "Refrax is an AI-powered refactoring agent for Java code. It communicates using the A2A protocol",
	}
	root.AddCommand(
		newStartCmd(),
	)
	return root
}
