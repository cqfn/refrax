package client

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/critic"
	"github.com/cqfn/refrax/internal/facilitator"
	"github.com/cqfn/refrax/internal/fixer"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
)

// RefraxClient represents a client used for refactoring projects.
type RefraxClient struct {
	provider string
	playbook string
	token    string
}

// NewRefraxClient creates a new instance of RefraxClient.
func NewRefraxClient(provider, token, playbook string) *RefraxClient {
	return &RefraxClient{
		provider: provider,
		token:    token,
		playbook: playbook,
	}
}

// Refactor initializes the refactoring process for the given project.
func Refactor(provider, token string, proj Project, stats bool, logger log.Logger, playbook string) (Project, error) {
	return NewRefraxClient(provider, token, playbook).Refactor(proj, stats, logger)
}

// Refactor performs refactoring on the given project using the RefraxClient.
func (c *RefraxClient) Refactor(proj Project, stats bool, logger log.Logger) (Project, error) {
	logger.Debug("starting refactoring for project %s", proj)
	classes, err := proj.Classes()
	if err != nil {
		return nil, fmt.Errorf("failed to get classes from project %s: %w", proj, err)
	}
	if len(classes) == 0 {
		return proj, fmt.Errorf("no java classes found in the project %s, add java files to the appropriate directory", proj)
	}
	logger.Debug("found %d classes in the project: %v", len(classes), classes)
	var ai brain.Brain
	mind := brain.New(c.provider, c.token, c.playbook)
	if stats {
		ai = brain.NewMetricBrain(mind, logger)
	} else {
		ai = mind
	}
	criticPort, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for critic: %w", err)
	}
	ctc := critic.NewCritic(ai, criticPort)

	fixerPort, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for fixer: %w", err)
	}
	fxr := fixer.NewFixer(ai, fixerPort)

	facilitatorPort, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for facilitator: %w", err)
	}
	fclttor := facilitator.NewFacilitator(ai, facilitatorPort, criticPort, fixerPort)

	facilitatorReady := make(chan struct{})
	criticReady := make(chan struct{})
	fixerReady := make(chan struct{})

	go func() {
		faerr := fclttor.Start(facilitatorReady)
		if faerr != nil && faerr != http.ErrServerClosed {
			panic(fmt.Sprintf("failed to start facilitator server: %v", faerr))
		}
	}()
	go func() {
		ferr := fxr.Start(fixerReady)
		if ferr != nil && ferr != http.ErrServerClosed {
			panic(fmt.Sprintf("failed to start fixer server: %v", ferr))
		}
	}()
	go func() {
		cerr := ctc.Start(criticReady)
		if cerr != nil && cerr != http.ErrServerClosed {
			panic(fmt.Sprintf("failed to start critic server: %v", cerr))
		}
	}()

	defer closeResource(ctc)
	defer closeResource(fclttor)
	defer closeResource(fxr)

	<-facilitatorReady
	<-criticReady
	<-fixerReady

	logger.Info("all servers are ready: facilitator %d, critic %d, fixer %d", facilitatorPort, criticPort, fixerPort)
	logger.Info("begin refactoring")
	facilitatorClient := protocol.NewCustomClient(fmt.Sprintf("http://localhost:%d", facilitatorPort))

	for _, class := range classes {
		logger.Debug("sending class %s for refactoring", class.Name())
		var resp *protocol.JSONRPCResponse
		resp, err = facilitatorClient.SendMessage(protocol.MessageSendParams{
			Message: protocol.NewMessageBuilder().
				MessageID("1").
				Part(protocol.NewText(fmt.Sprintf("Refactor the class '%s'", class.Name()))).
				Part(protocol.NewFileBytes([]byte(class.Content()))).
				Build(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send message for class %s: %w", class.Name(), err)
		}
		logger.Debug("received response for class %s: %s", class.Name(), resp)
		refactored := resp.Result.(protocol.Message).Parts[0].(*protocol.FilePart).File.(protocol.FileWithBytes).Bytes
		var decoded []byte
		decoded, err = base64.StdEncoding.DecodeString(refactored)
		if err != nil {
			return nil, fmt.Errorf("failed to decode refactored class %s: %w", class.Name(), err)
		}
		err = class.SetContent(clean(string(decoded)))
		if err != nil {
			return nil, fmt.Errorf("failed to set content for class %s: %w", class.Name(), err)
		}
	}
	logger.Info("refactoring is finished")

	if withStats, ok := ai.(*brain.MetricBrain); ok {
		withStats.PrintStats()
	}
	return proj, err
}

func clean(s string) string {
	return strings.TrimSpace(s)
}

func closeResource(resource io.Closer) {
	if cerr := resource.Close(); cerr != nil {
		panic(fmt.Sprintf("failed to close resource: %v", cerr))
	}
}
