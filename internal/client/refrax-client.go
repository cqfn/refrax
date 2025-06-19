package client

import (
	"io"

	"github.com/cqfn/refrax/internal/critic"
	"github.com/cqfn/refrax/internal/facilitator"
)

func Refactor(question string) (string, error) {
	facilitator, err := facilitator.NewFacilitator("none", 8080)
	if err != nil {
		return "", err
	}
	critic, err := critic.NewCritic("none", 8081)
	if err != nil {
		return "", err
	}
	fready := make(chan struct{})
	cready := make(chan struct{})
	go startServer(facilitator, fready, &err)
	go startCriticServer(critic, cready, &err)

	<-fready
	<-cready

	// Do some job

	defer closeResource(critic, &err)
	defer closeResource(facilitator, &err)

	// This function is a placeholder for the actual implementation.
	// In a real-world scenario, this would send the question to the server
	// and return the response.
	return "This is a mock response to the question: " + question, nil
}

func startServer(server *facilitator.Facilitator, ready chan struct{}, err *error) {
	if cerr := server.Start(ready); cerr != nil && *err == nil {
		*err = cerr
	}
}

func startCriticServer(server *critic.Critic, ready chan struct{}, err *error) {
	if cerr := server.Start(ready); cerr != nil && *err == nil {
		*err = cerr
	}
}

func closeResource(resource io.Closer, err *error) {
	if cerr := resource.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
