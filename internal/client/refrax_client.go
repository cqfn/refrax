package client

import (
	"encoding/base64"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/critic"
	"github.com/cqfn/refrax/internal/facilitator"
	"github.com/cqfn/refrax/internal/fixer"
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

func Refactor(provider string, token string, proj Project, stats bool, log log.Logger) (Project, error) {
	return NewRefraxClient(provider, token).Refactor(proj, stats, log)
}

func (c *RefraxClient) Refactor(proj Project, stats bool, log log.Logger) (Project, error) {
	cmd := exec.Command("aibolit", "check", "--filenames", "Foo.java")
	opportunities, _ := cmd.CombinedOutput()
    log.Debug("Identified refactoring opportunities with aibolit: \n%s", opportunities)
    // pass result to prompt
	log.Debug("starting refactoring for project %s", proj)
	classes, err := proj.Classes()
	if err != nil {
		return nil, fmt.Errorf("failed to get classes from project %s: %w", proj, err)
	}
	if len(classes) == 0 {
		return proj, fmt.Errorf("no java classes found in the project %s, add java files to the appropriate directory", proj)
	}
	log.Debug("found %d classes in the project: %v", len(classes), classes)

	var ai brain.Brain
	if stats {
		ai = brain.NewMetricBrain(brain.New(c.provider, c.token), log)
	} else {
		ai = brain.New(c.provider, c.token)
	}

	criticPort, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for critic: %w", err)
	}
	critic := critic.NewCritic(ai, criticPort, opportunities)

	fixerPort, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for fixer: %w", err)
	}
	fixer := fixer.NewFixer(ai, fixerPort)

	facilitatorPort, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for facilitator: %w", err)
	}
	facilitator := facilitator.NewFacilitator(ai, facilitatorPort, criticPort, fixerPort)

	facilitatorReady := make(chan struct{})
	criticReady := make(chan struct{})
	fixerReady := make(chan struct{})

	go startServer(facilitator, facilitatorReady, &err)
	go startCriticServer(critic, criticReady, &err)
	go startFixerServer(fixer, fixerReady, &err)
	defer closeResource(critic, &err)
	defer closeResource(facilitator, &err)
	defer closeResource(fixer, &err)

	<-facilitatorReady
	<-criticReady
	<-fixerReady

	log.Info("all servers are ready: facilitator %d, critic %d, fixer %d", facilitatorPort, criticPort, fixerPort)
	log.Info("begin refactoring")
	facilitatorClient := protocol.NewCustomClient(fmt.Sprintf("http://localhost:%d", facilitatorPort))

	for _, class := range classes {
		log.Debug("sending class %s for refactoring", class.Name())
		resp, err := facilitatorClient.SendMessage(protocol.MessageSendParams{
			Message: protocol.NewMessageBuilder().
				MessageID("1").
				Part(protocol.NewText(fmt.Sprintf("Refactor the class '%s'", class.Name()))).
				Part(protocol.NewFileBytes([]byte(class.Content()))).
				Build(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send message for class %s: %w", class.Name(), err)
		}
		log.Debug("received response for class %s: %s", class.Name(), resp)
		refactored := resp.Result.(protocol.Message).Parts[0].(*protocol.FilePart).File.(protocol.FileWithBytes).Bytes
		decoded, err := base64.StdEncoding.DecodeString(refactored)
		if err != nil {
			return nil, fmt.Errorf("failed to decode refactored class %s: %w", class.Name(), err)
		}
		err = class.SetContent(clean(string(decoded)))
		if err != nil {
			return nil, fmt.Errorf("failed to set content for class %s: %w", class.Name(), err)
		}
	}
	log.Info("refactoring is finished")
	if withStats, ok := ai.(*brain.MetricBrain); ok {
		withStats.PrintStats()
	}
	return proj, err
}

func clean(s string) string {
	return strings.TrimSpace(s)
}

func startFixerServer(fixer *fixer.Fixer, ready chan struct{}, err *error) {
	if cerr := fixer.Start(ready); cerr != nil && *err == nil {
		*err = cerr
	}
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
