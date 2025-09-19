package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrentVersion(t *testing.T) {
	datetime = "2024-01-01T00:00:00Z"
	assert.Equal(t, "dev built at "+datetime, Version())
}
