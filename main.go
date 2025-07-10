// Refrax entry point.
package main

import (
	"fmt"
	"os"

	"github.com/cqfn/refrax/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
