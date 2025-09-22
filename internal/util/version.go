package util

import (
	"runtime/debug"
	"strings"
)

var (
	// Use the following command to change this value:
	// go build -ldflags "-X github.com/cqfn/refrax/internal/util.version=<version>"
	version string = ""

	// Use the following command to change this value:
	// go build -ldflags "-X github.com/cqfn/refrax/internal/util.commit=<commit-hash>"
	commit string = ""

	// Use the following command to change this value:
	// go build -ldflags "-X github.com/cqfn/refrax/internal/util.datetime=<date-time>"
	datetime string = ""
)

func Version() string {
	info := info()
	parts := []string{ver(info)}
	rev := revision(info)
	if rev != "" {
		parts = append(parts, "commit", rev)
	}
	time := timestamp(info)
	if time != "" {
		parts = append(parts, "built at", time)
	}
	return strings.Join(parts, " ")
}

func ver(info map[string]string) string {
	if version != "" {
		return version
	}
	if iver, ok := info["version"]; ok {
		return iver
	}
	return "dev"
}

func revision(info map[string]string) string {
	if commit != "" {
		return commit
	}
	if rev, ok := info["vcs.revision"]; ok {
		return rev
	}
	return ""
}

func timestamp(info map[string]string) string {
	if datetime != "" {
		return datetime
	}
	if time, ok := info["vcs.time"]; ok {
		return time
	}
	return ""
}

func info() map[string]string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("failed to read build info")
	}
	settings := make(map[string]string)
	for _, s := range info.Settings {
		settings[s.Key] = s.Value
	}
	vcsver := info.Main.Version
	if vcsver != "" && vcsver != "(devel)" {
		settings["version"] = vcsver
	}
	return settings
}
