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
	"github.com/cqfn/refrax/internal/prompts"
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
	log.Debug("Starting refactoring for project %s", proj)
	classes, err := proj.Classes()
	if err != nil {
		return nil, fmt.Errorf("failed to get classes from project %s: %w", proj, err)
	}
	if len(classes) == 0 {
		return proj, fmt.Errorf("no java classes found in the project %s, add java files to the appropriate directory", proj)
	}
	log.Debug("Found %d classes in the project: %v", len(classes), classes)

	criticStats := &stats.Stats{Name: "critic"}
	criticSystemPrompt := prompts.System{
		AgentName:      "critic",
		ProjectContext: "you are part of a team working on a Java project. Your role is to review Java classes and provide constructive feedback to improve code quality, maintainability, and adherence to best practices.",
		Capabilities: []string{
			"Analyze Java code for potential improvements",
			"Identify code smells and suggest refactorings",
			"Provide feedback on code structure and design patterns",
			"Suggest improvements without altering functionality",
		},
		Constraints: []string{
			"You cannot change the functionality of the code",
			"You cannot suggest changes that require moving code between files",
			"You cannot suggest renamimg classes or methods",
			"You cannot suggest removing JavaDoc comments",
		},
	}
	token, err := token(c.params)
	if err != nil {
		return nil, fmt.Errorf("failed to find token: %w", err)
	}
	criticBrain, err := mind(c.params, token, &criticSystemPrompt, criticStats)
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
	fixerSystemPrompt := prompts.System{
		AgentName:      "fixer",
		ProjectContext: "you are part of a team working on a Java project. Your role is to fix Java classes based on the feedback provided by the Critic, ensuring that the code quality and maintainability are improved without altering the original functionality.",
		Capabilities: []string{
			"Apply suggested improvements to Java code",
			"Refactor code to enhance readability and maintainability",
		},
		Constraints: []string{
			"You cannot change the functionality of the code",
			"You cannot change the code that require moving code between files",
			"You cannot rename classes or methods",
			"You cannot remove JavaDoc comments",
		},
	}
	fixerBrain, err := mind(c.params, token, &fixerSystemPrompt, fixerStats)
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
	reviewerSystemPrompt := prompts.System{
		AgentName:      "reviewer",
		ProjectContext: "you are part of a team working on a Java project. Your role is to review the refactored Java classes to ensure that the applied changes align with the original suggestions provided by the Critic and that the code quality has been improved without altering the original functionality.",
		Capabilities: []string{
			"Run build and test commands to validate code changes",
			"Provide feedback on the success or failure of the build and tests",
			"Suggest further improvements based on build and test results",
		},
		Constraints: []string{
			"You cannot suggest changes that require moving code between files",
			"You cannot suggest adding another dependencies",
		},
	}
	reviewerBrain, err := mind(c.params, token, &reviewerSystemPrompt, reviewerStats)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI instance for reviewer: %w", err)
	}
	reviewerPort, err := util.FreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find free port for reviewer: %w", err)
	}
	rvwr := reviewer.NewReviewer(reviewerBrain, reviewerPort, c.params.Checks...)
	rvwr.Handler(countStats(reviewerStats))

	facilitatorStats := &stats.Stats{Name: "facilitator"}
	facilitatorSystemPrompt := prompts.System{
		AgentName:      "facilitator",
		ProjectContext: "you are part of a team working on a Java project. Your role is to facilitate the refactoring process by coordinating between the Critic, Fixer, and Reviewer agents to ensure that Java classes are effectively improved while maintaining their original functionality.",
		Capabilities: []string{
			"Understand the most important suggestions from the Critic",
			"Group and prioritize suggestions for the Fixer",
		},
		Constraints: []string{
			"You cannot change suggestions",
		},
	}
	facilitatorBrain, err := mind(c.params, token, &facilitatorSystemPrompt, facilitatorStats)
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
			panic(fmt.Sprintf("Failed to start facilitator server: %v", faerr))
		}
	}()
	go func() {
		ferr := fxr.ListenAndServe()
		if ferr != nil && ferr != http.ErrServerClosed {
			panic(fmt.Sprintf("Failed to start fixer server: %v", ferr))
		}
	}()
	go func() {
		cerr := ctc.ListenAndServe()
		if cerr != nil && cerr != http.ErrServerClosed {
			panic(fmt.Sprintf("Failed to start critic server: %v", cerr))
		}
	}()
	go func() {
		rerr := rvwr.ListenAndServe()
		if rerr != nil && rerr != http.ErrServerClosed {
			panic(fmt.Sprintf("Failed to start reviewer server: %v", rerr))
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

	log.Info("All servers are ready: facilitator %d, critic %d, fixer %d, reviewer %d", facilitatorPort, criticPort, fixerPort, reviewerPort)
	log.Info("Begin refactoring for project %s with %d classes", proj, len(classes))
	ch := make(chan refactoring, len(classes))
	go refactor(fclttor, proj, c.params.MaxSize, ch)
	for range len(classes) {
		res := <-ch
		if res.class != nil && res.content != "" {
			log.Info("Received refactored class: %s, content length: %d", res.class.Name(), len(res.content))
		}
	}
	log.Info("Refactoring is finished")
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
	log.Debug("Refactoring project %q", p)
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
	job := domain.Job{
		Descr: &domain.Description{
			Text: "refactor the project",
			Meta: map[string]any{
				"max-size": fmt.Sprintf("%d", size),
			},
		},
		Classes: all,
	}
	artifacts, err := f.Refactor(&job)
	if err != nil {
		log.Error("Failed to refactor project %s: %v", p, err)
		ch <- refactoring{err: fmt.Errorf("failed to refactor project %s: %w", p, err)}
		close(ch)
		return
	}
	refactored := artifacts.Classes
	log.Info("Refactored %d classes in project %s", len(refactored), p)
	for _, c := range refactored {
		log.Debug("Received refactored class: ", c)
		ch <- refactoring{class: before[c.Name()], content: c.Content(), err: nil}
	}
	close(ch)
}

type shudownable interface {
	Shutdown() error
}

func shutdown(s shudownable) {
	if cerr := s.Shutdown(); cerr != nil {
		panic(fmt.Sprintf("Failed to close resource: %v", cerr))
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
			log.Info("Using csv file for statistics output")
			output := p.Soutput
			if output == "" {
				output = "stats.csv"
			}
			swriter = stats.NewCSVWriter(output)
		} else {
			log.Info("Using stdout format for statistics output")
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

func mind(p Params, token string, system *prompts.System, s *stats.Stats) (brain.Brain, error) {
	ai, err := brain.New(p.Provider, token, system.String(), p.Playbook)
	if p.Stats {
		ai = brain.NewMetricBrain(ai, s)
	}
	return ai, err
}

func token(p Params) (string, error) {
	log.Debug("Refactoring provider: %s", p.Provider)
	log.Debug("Project path to refactor: %s", p.Input)
	var token string
	if p.Token != "" {
		token = p.Token
	} else {
		log.Info("Token not provided, trying to find token in .env file")
		token = env.Token(".env", p.Provider)
	}
	if token == "" {
		return "", fmt.Errorf("token not found, please provide it via --token flag or in .env file")
	}
	log.Debug("Using provided token: %s...", mask(token))
	return token, nil
}

func proj(params Params) (domain.Project, error) {
	if params.MockProject {
		log.Debug("Using mock project")
		return domain.NewMock(), nil
	}
	input := domain.NewFilesystem(params.Input)
	output := params.Output
	if output != "" {
		log.Debug("Copy project to %q", output)
		return domain.NewMirrorProject(input, output)
	}
	log.Debug("No output path provided, changing project in place %q", params.Input)
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
