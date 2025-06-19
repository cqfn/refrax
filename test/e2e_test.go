package test

import (
	"bytes"
	"io"
	"testing"

	"github.com/cqfn/refrax/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndToEnd_Agents_FromCLI(t *testing.T) {
	output := &bytes.Buffer{}
	cmd := cmd.NewRootCmd(output, io.Discard)
	cmd.SetArgs([]string{"refactor", "--ai=none"})

	err := cmd.Execute()

	require.NoError(t, err, "Expected command to execute without error")
	assert.Contains(t, output.String(), "provider: none", "expect no AI provider to be used in output")
}
