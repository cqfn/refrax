package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/critic"
	"github.com/cqfn/refrax/internal/env"
	"github.com/cqfn/refrax/internal/facilitator"
	"github.com/cqfn/refrax/internal/fixer"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/stats"
)

// RefraxClient represents a client used for refactoring projects.
type RefraxClient struct {
	params Params
}

// NewRefraxClient creates a new instance of RefraxClient.
func NewRefraxClient(params *Params) *RefraxClient {
	initLogger(params)
	return &RefraxClient{
		params: *params,
	}
}

// Refactor initializes the refactoring process for the given project.
func Refactor(params *Params) (Project, error) {
	proj, err := project(*params)
	if err != nil {
		return nil, fmt.Errorf("failed to create project from params: %w", err)
	}
	return NewRefraxClient(params).Refactor(proj)
}

// Refactor performs refactoring on the given project using the RefraxClient.
func (c *RefraxClient) Refactor(proj Project) (Project, error) {
	log.Debug("starting refactoring for project %s", proj)
	classes, err := proj.Classes()
	if err != nil {
		return nil, fmt.Errorf("failed to get classes from project %s: %w", proj, err)
	}
	if len(classes) == 0 {
		return proj, fmt.Errorf("no java classes found in the project %s, add java files to the appropriate directory", proj)
	}
	log.Debug("found %d classes in the project: %v", len(classes), classes)
	counter := &stats.Stats{}
	ai := mind(c.params, counter)

	criticPort, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for critic: %w", err)
	}
	ctc := critic.NewCritic(ai, criticPort)
	ctc.Handler(countStats(counter))

	fixerPort, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for fixer: %w", err)
	}
	fxr := fixer.NewFixer(ai, fixerPort)
	fxr.Handler(countStats(counter))

	facilitatorPort, err := protocol.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for facilitator: %w", err)
	}
	fclttor := facilitator.NewFacilitator(ai, facilitatorPort, criticPort, fixerPort)
	fclttor.Handler(countStats(counter))

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

	log.Info("all servers are ready: facilitator %d, critic %d, fixer %d", facilitatorPort, criticPort, fixerPort)
	log.Info("begin refactoring")
	facilitatorClient := protocol.NewCustomClient(fmt.Sprintf("http://localhost:%d", facilitatorPort))

	for _, class := range classes {
		log.Debug("sending class %s for refactoring", class.Name())
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
		log.Debug("received response for class %s: %s", class.Name(), resp)
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
	log.Info("refactoring is finished")
	err = printStats(c.params, counter)
	if err != nil {
		return nil, fmt.Errorf("failed to print statistics: %w", err)
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

func initLogger(params *Params) {
	if params.Debug {
		log.Set(log.NewZerolog(params.Log, "debug"))
	} else {
		log.Set(log.NewZerolog(params.Log, "info"))
	}
}

func printStats(p Params, s *stats.Stats) error {
	if p.Stats {
		var swriter stats.Writer
		if p.Format == "csv" {
			log.Info("using csv file for statistics output")
			output := p.Soutput
			if output == "" {
				output = "stats.csv"
			}
			swriter = stats.NewCSVWriter(output)
		} else {
			log.Info("using stdout format for statistics output")
			swriter = stats.NewStdWriter(log.Default())
		}
		return swriter.Print(s)
	}
	return nil
}

func mind(p Params, s *stats.Stats) brain.Brain {
	var ai brain.Brain
	ai = brain.New(p.Provider, token(p), p.Playbook)
	if p.Stats {
		ai = brain.NewMetricBrain(ai, s)
	}
	return ai
}

func token(p Params) string {
	log.Debug("refactoring provider: %s", p.Provider)
	log.Debug("project path to refactor: %s", p.Input)
	var token string
	if p.Token != "" {
		token = p.Token
	} else {
		log.Info("token not provided, trying to find token in .env file")
		token = env.Token(".env", p.Provider)
	}
	log.Debug("using provided token: %s...", mask(token))
	return token
}

func project(params Params) (Project, error) {
	if params.Mock {
		log.Debug("using mock project")
		return NewMockProject(), nil
	}
	input := NewFilesystemProject(params.Input)
	output := params.Output
	if output != "" {
		log.Debug("copy project to %q", output)
		return NewMirrorProject(input, output)
	}
	log.Debug("no output path provided, changing project in place %q", params.Input)
	return input, nil
}

func mask(token string) string {
	n := len(token)
	if n == 0 {
		return ""
	}
	visible := min(n, 3)
	return token[:visible] + strings.Repeat("*", n-visible)
}

func countStats(s *stats.Stats) protocol.Handler {
	return func(next protocol.Handler, r *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
		start := time.Now()
		resp, err := next(nil, r)
		if err != nil {
			return nil, fmt.Errorf("failed to process request: %w", err)
		}
		duration := time.Since(start)
		jsonresp, err := json.Marshal(resp)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		jsonreq, err := json.Marshal(r)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqt, err := stats.Tokens(string(jsonreq))
		if err != nil {
			return nil, fmt.Errorf("failed to count tokens for request: %w", err)
		}
		respt, err := stats.Tokens(string(jsonresp))
		if err != nil {
			return nil, fmt.Errorf("failed to count tokens for response: %w", err)
		}
		s.A2AReq(duration, reqt, respt, len(jsonreq), len(jsonresp))
		return resp, err
	}
}
