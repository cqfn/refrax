package client

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/cqfn/refrax/internal/critic"
	"github.com/cqfn/refrax/internal/facilitator"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

type RefraxClient struct {
	client protocol.Client
}

func NewRefraxClient() *RefraxClient {
	return &RefraxClient{
		client: protocol.NewCustomClient("http://localhost:8080"),
	}
}

func Refactor(proj Project) (Project, error) {
	return NewRefraxClient().Refactor(proj)
}

func (c *RefraxClient) Refactor(proj Project) (Project, error) {
	log.Debug("starting refactoring for project %s", proj)
	facilitator, err := facilitator.NewFacilitator("none", 8080)
	if err != nil {
		return nil, err
	}
	critic, err := critic.NewCritic("none", 8081)
	if err != nil {
		return nil, err
	}
	fready := make(chan struct{})
	cready := make(chan struct{})
	go startServer(facilitator, fready, &err)
	go startCriticServer(critic, cready, &err)

	<-fready
	<-cready

	all, err := proj.Classes()
	if err != nil {
		return nil, err
	}
	log.Debug("found %d classes in the project: %v", len(all), all)
	for _, class := range all {
		log.Debug("sending class %s for refactoring", class.Name())
		resp, err := c.client.SendMessage(protocol.MessageSendParams{
			Message: protocol.Message{
				MessageID: "1",
				Parts: protocol.Parts{
					&protocol.TextPart{
						Text: fmt.Sprintf("Refactor the class '%s'", class.Name()),
						Kind: protocol.PartKindText,
					}, &protocol.FilePart{
						Kind: protocol.PartKindFile,
						File: protocol.FileWithBytes{
							Bytes: base64.StdEncoding.EncodeToString([]byte(class.Content())),
						},
					},
				},
			}})
		if err != nil {
			return nil, fmt.Errorf("failed to send message for class %s: %w", class.Name(), err)
		}
		refactored := resp.Result.(protocol.Message).Parts[0].(*protocol.FilePart).File.(protocol.FileWithBytes).Bytes
		decoded, err := base64.StdEncoding.DecodeString(refactored)
		if err != nil {
			return nil, fmt.Errorf("failed to decode refactored class %s: %w", class.Name(), err)
		}
		class.SetContent(string(decoded))
		log.Info("client refactored the class %s: %s", class.Name(), class.Content())
	}
	defer closeResource(critic, &err)
	defer closeResource(facilitator, &err)
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
