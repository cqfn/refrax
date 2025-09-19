package util

import (
	"strings"
)

const (
	defaultVersion  = "dev"
	defaultCommit   = "none"
	defaultDatetime = "none"
)

var (
	// Use the following command to change this value:
	// go build -ldflags "-X github.com/cqfn/refrax/internal/util.version=<version>"
	version string = defaultVersion

	// Use the following command to change this value:
	// go build -ldflags "-X github.com/cqfn/refrax/internal/util.commit=<commit-hash>"
	commit string = defaultCommit

	// Use the following command to change this value:
	// go build -ldflags "-X github.com/cqfn/refrax/internal/util.datetime=<date-time>"
	datetime string = defaultDatetime
)

func Version() string {
	parts := []string{version}
	if commit != defaultCommit {
		parts = append(parts, "commit", commit)
	}
	if datetime != defaultDatetime {
		parts = append(parts, "built at", datetime)
	}
	return strings.Join(parts, " ")
}
