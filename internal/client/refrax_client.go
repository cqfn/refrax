package client

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/critic"
	"github.com/cqfn/refrax/internal/facilitator"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

type RefraxClient struct {
	provider string
	token    string
}

func NewRefraxClient(provider string, token string) *RefraxClient {
	return &RefraxClient{
		provider: provider,
		token:    token,
	}
}

func Refactor(provider string, token string, proj Project) (Project, error) {
	return NewRefraxClient(provider, token).Refactor(proj)
}

func (c *RefraxClient) Refactor(proj Project) (Project, error) {
	log.Debug("starting refactoring for project %s", proj)
	fport, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for refactoring: %w", err)
	}
	facilitator, err := facilitator.NewFacilitator(brain.New(c.provider, c.token), fport)
	if err != nil {
		return nil, err
	}
	cport, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for refactoring: %w", err)
	}
	critic, err := critic.NewCritic(c.provider, cport)
	if err != nil {
		return nil, err
	}
	fready := make(chan struct{})
	cready := make(chan struct{})
	go startServer(facilitator, fready, &err)
	go startCriticServer(critic, cready, &err)
	defer closeResource(critic, &err)
	defer closeResource(facilitator, &err)

	<-fready
	<-cready

	client := protocol.NewCustomClient(fmt.Sprintf("http://localhost:%d", fport))

	all, err := proj.Classes()
	if err != nil {
		return nil, err
	}
	log.Debug("found %d classes in the project: %v", len(all), all)
	for _, class := range all {
		log.Debug("sending class %s for refactoring", class.Name())
		resp, err := client.SendMessage(protocol.MessageSendParams{
			Message: protocol.NewMessageBuilder().
				MessageID("1").
				Part(protocol.NewText(fmt.Sprintf("Refactor the class '%s'", class.Name()))).
				Part(protocol.NewFileBytes([]byte(class.Content()))).
				Build(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send message for class %s: %w", class.Name(), err)
		}
		refactored := resp.Result.(protocol.Message).Parts[0].(*protocol.FilePart).File.(protocol.FileWithBytes).Bytes
		decoded, err := base64.StdEncoding.DecodeString(refactored)
		if err != nil {
			return nil, fmt.Errorf("failed to decode refactored class %s: %w", class.Name(), err)
		}
		err = class.SetContent(string(decoded))
		if err != nil {
			return nil, fmt.Errorf("failed to set content for class %s: %w", class.Name(), err)
		}
		log.Info("client refactored the class %s: %s", class.Name(), class.Content())
	}
	return proj, err
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
