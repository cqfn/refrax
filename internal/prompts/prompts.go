// Package prompts is for the prompts.
package prompts

import (
	"embed"
	"strings"
	"text/template"
)

//go:embed *.tmpl */*.tmpl
var files embed.FS

type System struct {
	AgentName      string
	ProjectContext string
	Constraints    []string
	Capabilities   []string
}

func (s *System) String() string {
	tmpl, err := template.ParseFS(files, "system.md.tmpl")
	if err != nil {
		panic(err)
	}
	var result strings.Builder
	err = tmpl.Execute(&result, s)
	if err != nil {
		panic(err)
	}
	return result.String()
}

type User struct {
	Data any
	Name string
}

func (u *User) String() string {
	tmpl, err := template.ParseFS(files, u.Name)
	if err != nil {
		panic(err)
	}
	var result strings.Builder
	err = tmpl.Execute(&result, u.Data)
	if err != nil {
		panic(err)
	}
	return result.String()
}
