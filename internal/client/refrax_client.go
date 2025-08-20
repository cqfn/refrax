package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/critic"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/env"
	"github.com/cqfn/refrax/internal/facilitator"
	"github.com/cqfn/refrax/internal/fixer"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/protocol"
	"github.com/cqfn/refrax/internal/reviewer"
	"github.com/cqfn/refrax/internal/stats"
	"github.com/cqfn/refrax/internal/util"
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
func Refactor(params *Params) (domain.Project, error) {
	proj, err := proj(*params)
	if err != nil {
		return nil, fmt.Errorf("failed to create project from params: %w", err)
	}
	return NewRefraxClient(params).Refactor(proj)
}

// Refactor performs refactoring on the given project using the RefraxClient.
func (c *RefraxClient) Refactor(proj domain.Project) (domain.Project, error) {
	log.Debug("starting refactoring for project %s", proj)
	classes, err := proj.Classes()
	if err != nil {
		return nil, fmt.Errorf("failed to get classes from project %s: %w", proj, err)
	}
	if len(classes) == 0 {
		return proj, fmt.Errorf("no java classes found in the project %s, add java files to the appropriate directory", proj)
	}
	log.Debug("found %d classes in the project: %v", len(classes), classes)

	criticStats := &stats.Stats{Name: "critic"}
	criticBrain, err := mind(c.params, criticStats)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI instance: %w", err)
	}
	criticPort, err := util.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for critic: %w", err)
	}
	ctc := critic.NewCritic(criticBrain, criticPort)
	ctc.Handler(countStats(criticStats))

	fixerStats := &stats.Stats{Name: "fixer"}
	fixerBrain, err := mind(c.params, fixerStats)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI instance: %w", err)
	}
	fixerPort, err := util.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for fixer: %w", err)
	}
	fxr := fixer.NewFixer(fixerBrain, fixerPort)
	fxr.Handler(countStats(fixerStats))

	reviewerStats := &stats.Stats{Name: "reviewer"}
	reviewerBrain, err := mind(c.params, reviewerStats)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI instance for reviewer: %w", err)
	}
	reviewerPort, err := util.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for reviewer: %w", err)
	}
	rvwr := reviewer.NewReviewer(reviewerBrain, reviewerPort, "mvn clean test", "mvn qulice:check -Pqulice")
	rvwr.Handler(countStats(reviewerStats))

	facilitatorStats := &stats.Stats{Name: "facilitator"}
	facilitatorBrain, err := mind(c.params, facilitatorStats)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI instance: %w", err)
	}
	facilitatorPort, err := util.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for facilitator: %w", err)
	}
	fclttor := facilitator.NewFacilitator(facilitatorBrain, ctc, fxr, rvwr, facilitatorPort)
	fclttor.Handler(countStats(facilitatorStats))

	go func() {
		faerr := fclttor.ListenAndServe()
		if faerr != nil && faerr != http.ErrServerClosed {
			panic(fmt.Sprintf("failed to start facilitator server: %v", faerr))
		}
	}()
	go func() {
		ferr := fxr.ListenAndServe()
		if ferr != nil && ferr != http.ErrServerClosed {
			panic(fmt.Sprintf("failed to start fixer server: %v", ferr))
		}
	}()
	go func() {
		cerr := ctc.ListenAndServe()
		if cerr != nil && cerr != http.ErrServerClosed {
			panic(fmt.Sprintf("failed to start critic server: %v", cerr))
		}
	}()
	go func() {
		rerr := rvwr.ListenAndServe()
		if rerr != nil && rerr != http.ErrServerClosed {
			panic(fmt.Sprintf("failed to start reviewer server: %v", rerr))
		}
	}()

	defer shutdown(ctc)
	defer shutdown(fclttor)
	defer shutdown(fxr)
	defer shutdown(rvwr)

	<-ctc.Ready()
	<-fclttor.Ready()
	<-fxr.Ready()
	<-rvwr.Ready()

	log.Info("all servers are ready: facilitator %d, critic %d, fixer %d, reviewer %d", facilitatorPort, criticPort, fixerPort, reviewerPort)
	log.Info("begin refactoring")
	ch := make(chan refactoring, len(classes))
	go refactor(fclttor, proj, c.params.MaxSize, ch)
	for range len(classes) {
		res := <-ch
		if res.class != nil && res.content != "" {
			log.Info("received refactored class: %s, content length: %d", res.class.Name(), len(res.content))
		}
	}
	log.Info("refactoring is finished")
	err = printStats(c.params, criticStats, fixerStats, facilitatorStats)
	if err != nil {
		return nil, fmt.Errorf("failed to print statistics: %w", err)
	}
	return proj, err
}

type refactoring struct {
	class   domain.Class
	content string
	err     error
}

func refactor(f domain.Facilitator, p domain.Project, size int, ch chan<- refactoring) {
	log.Debug("refactoring project %q", p)
	all, err := p.Classes()
	if err != nil {
		ch <- refactoring{err: fmt.Errorf("failed to get classes from project %s: %w", p, err)}
		close(ch)
		return
	}
	before := make(map[string]domain.Class)
	for _, c := range all {
		before[c.Name()] = c
	}
	task := domain.NewTask("refactor the project", all, map[string]any{"max-size": fmt.Sprintf("%d", size)})
	refactored, err := f.Refactor(task)
	if err != nil {
		ch <- refactoring{err: fmt.Errorf("failed to refactor project %s: %w", p, err)}
		close(ch)
		return
	}
	log.Info("refactored %d classes in project %s", len(refactored), p)
	for _, c := range refactored {
		log.Debug("rececived refactored class: ", c)
		ch <- refactoring{class: before[c.Name()], content: c.Content(), err: nil}
	}
	close(ch)
}

type shudownable interface {
	Shutdown() error
}

func shutdown(s shudownable) {
	if cerr := s.Shutdown(); cerr != nil {
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

func printStats(p Params, s ...*stats.Stats) error {
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
		var res []*stats.Stats
		total := &stats.Stats{}
		for _, st := range s {
			res = append(res, st)
			total = total.Add(st)
		}
		total.Name = "total"
		res = append(res, total)
		return swriter.Print(res...)
	}
	return nil
}

func mind(p Params, s *stats.Stats) (brain.Brain, error) {
	ai, err := brain.New(p.Provider, token(p), p.Playbook)
	if p.Stats {
		ai = brain.NewMetricBrain(ai, s)
	}
	return ai, err
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

func proj(params Params) (domain.Project, error) {
	if params.MockProject {
		log.Debug("using mock project")
		return domain.NewMock(), nil
	}
	input := domain.NewFilesystem(params.Input)
	output := params.Output
	if output != "" {
		log.Debug("copy project to %q", output)
		return domain.NewMirrorProject(input, output)
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
