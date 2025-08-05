package stats

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkoukk/tiktoken-go"
	"github.com/pkoukk/tiktoken-go-loader"
)

// Stats is a struct that contains interaction statistics between components.
type Stats struct {
	// name is a string that represents the name of the statistics.
	Name string

	// mu is a mutex to protect concurrent write access to stats.
	mu sync.Mutex

	// llmreq is a slice of time.Duration that stores the duration of each request made to the LLM.
	llmreq []time.Duration

	// llmreqtokens is an int that stores the number of tokens used in LLM requests.
	llmreqtokens int

	// llmresptokens is an int that stores the number of bytes received in LLM responses.
	llmreqbytes int

	// llmresptokens is an int that stores the number of tokens received in LLM responses.
	llmresptokens int

	// llmrespbytes is an int that stores the number of bytes received in LLM responses.
	llmrespbytes int

	// a2areqs is a slice of time.Duration that stores the duration of each request made to the A2A service.
	a2areqs []time.Duration

	// a2areqtokens is an int that stores the number of tokens used in A2A requests.
	a2areqtokens int

	// a2areqbytes is an int that stores the number of bytes used in A2A requests.
	a2areqbytes int

	// a2aresptokens is an int that stores the number of tokens received in A2A responses.
	a2aresptokens int

	// a2arespbytes is an int that stores the number of bytes received in A2A responses.
	a2arespbytes int
}

// AverageA2ARespBytes calculates the average number of A2A response bytes.
func (s *Stats) AverageA2ARespBytes() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.a2areqs) == 0 {
		return 0
	}
	return float64(s.a2arespbytes) / float64(len(s.a2areqs))
}

// AverageA2AReqBytes calculates the average number of A2A request bytes.
func (s *Stats) AverageA2AReqBytes() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.a2areqs) == 0 {
		return 0
	}
	return float64(s.a2areqbytes) / float64(len(s.a2areqs))
}

// AverageA2ARespTokens calculates the average number of A2A response tokens.
func (s *Stats) AverageA2ARespTokens() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.a2areqs) == 0 {
		return 0
	}
	return float64(s.a2aresptokens) / float64(len(s.a2areqs))
}

// AverageA2AReqTokens calculates the average number of A2A request tokens.
func (s *Stats) AverageA2AReqTokens() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.a2areqs) == 0 {
		return 0
	}
	return float64(s.a2areqtokens) / float64(len(s.a2areqs))
}

// AverageA2AReqDuration calculates the average duration of A2A requests.
func (s *Stats) AverageA2AReqDuration() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.a2areqs) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range s.a2areqs {
		total += d
	}
	return total / time.Duration(len(s.a2areqs))
}

// TotalLLMRequests returns a copy of the total LLM requests durations.
func (s *Stats) TotalLLMRequests() []time.Duration {
	duplicate := make([]time.Duration, len(s.llmreq))
	copy(duplicate, s.llmreq)
	return duplicate
}

// TotalA2ARespBytes calculates the total A2A response bytes.
func (s *Stats) TotalA2ARespBytes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.a2arespbytes
}

// TotalA2AReqBytes calculates the total A2A request bytes.
func (s *Stats) TotalA2AReqBytes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.a2areqbytes
}

// TotalA2ABytes calculates the total A2A bytes (request + response).
func (s *Stats) TotalA2ABytes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.a2areqbytes + s.a2arespbytes
}

// TotalA2ARespTokens calculates the total A2A response tokens.
func (s *Stats) TotalA2ARespTokens() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.a2aresptokens
}

// TotalA2AReqTokens calculates the total A2A request tokens.
func (s *Stats) TotalA2AReqTokens() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.a2areqtokens
}

// TotalA2ATokens calculates the total A2A tokens (request + response).
func (s *Stats) TotalA2ATokens() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.a2areqtokens + s.a2aresptokens
}

// TotalA2AReqDuration calculates the total duration of A2A requests.
func (s *Stats) TotalA2AReqDuration() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	var total time.Duration
	for _, d := range s.a2areqs {
		total += d
	}
	return total
}

// A2AMessages returns the number of A2A requests.
func (s *Stats) A2AMessages() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.a2areqs)
}

// AverageLLMRespBytes calculates the average number of LLM response bytes.
func (s *Stats) AverageLLMRespBytes() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.llmreq) == 0 {
		return 0
	}
	return float64(s.llmrespbytes) / float64(len(s.llmreq))
}

// AverageLLMReqBytes calculates the average number of LLM request bytes.
func (s *Stats) AverageLLMReqBytes() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.llmreq) == 0 {
		return 0
	}
	return float64(s.llmreqbytes) / float64(len(s.llmreq))
}

// AverageLLMRespTokens calculates the average number of LLM response tokens.
func (s *Stats) AverageLLMRespTokens() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.llmreq) == 0 {
		return 0
	}
	return float64(s.llmresptokens) / float64(len(s.llmreq))
}

// AverageLLMReqTokens calculates the average number of LLM request tokens.
func (s *Stats) AverageLLMReqTokens() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.llmreq) == 0 {
		return 0
	}
	return float64(s.llmreqtokens) / float64(len(s.llmreq))
}

// AverageLLMReqDuration calculates the average duration of LLM requests.
func (s *Stats) AverageLLMReqDuration() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.llmreq) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range s.llmreq {
		total += d
	}
	return total / time.Duration(len(s.llmreq))
}

// TotalLLMRespBytes calculates the total LLM response bytes.
func (s *Stats) TotalLLMRespBytes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.llmrespbytes
}

// TotalLLMReqBytes calculates the total LLM request bytes.
func (s *Stats) TotalLLMReqBytes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.llmreqbytes
}

// TotalLLMBytes calculates the total LLM bytes (request + response).
func (s *Stats) TotalLLMBytes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.llmreqbytes + s.llmrespbytes
}

// TotalLLMRespTokens calculates the total LLM response tokens.
func (s *Stats) TotalLLMRespTokens() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.llmresptokens
}

// TotalLLMReqTokens calculates the total LLM request tokens.
func (s *Stats) TotalLLMReqTokens() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.llmreqtokens
}

// TotalLLMTokens calculates the total LLM tokens (request + response).
func (s *Stats) TotalLLMTokens() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.llmreqtokens + s.llmresptokens
}

// TotalLLMReqDuration calculates the total duration of LLM requests.
func (s *Stats) TotalLLMReqDuration() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	var total time.Duration
	for _, d := range s.llmreq {
		total += d
	}
	return total
}

// LLMReq records a request statistics made to the LLM.
func (s *Stats) LLMReq(duration time.Duration, reqt, respt, reqb, respb int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.llmreq = append(s.llmreq, duration)
	s.llmreqtokens += reqt
	s.llmresptokens += respt
	s.llmreqbytes += reqb
	s.llmrespbytes += respb
}

// A2AReq records a request statistics made to the A2A service.
func (s *Stats) A2AReq(duration time.Duration, reqt, respt, reqb, respb int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.a2areqs = append(s.a2areqs, duration)
	s.a2areqtokens += reqt
	s.a2aresptokens += respt
	s.a2areqbytes += reqb
	s.a2arespbytes += respb
}

// Add combines the current Stats instance with another and returns a new Stats instance.
// It does not mutate either of the original Stats instances.
func (s *Stats) Add(other *Stats) *Stats {
	s.mu.Lock()
	defer s.mu.Unlock()
	other.mu.Lock()
	defer other.mu.Unlock()
	combined := &Stats{
		Name:          fmt.Sprintf("%s + %s", s.Name, other.Name),
		llmreq:        append([]time.Duration{}, s.llmreq...),
		llmreqtokens:  s.llmreqtokens,
		llmresptokens: s.llmresptokens,
		llmreqbytes:   s.llmreqbytes,
		llmrespbytes:  s.llmrespbytes,
		a2areqs:       append([]time.Duration{}, s.a2areqs...),
		a2areqtokens:  s.a2areqtokens,
		a2aresptokens: s.a2aresptokens,
		a2areqbytes:   s.a2areqbytes,
		a2arespbytes:  s.a2arespbytes,
	}
	// Append other durations
	combined.llmreq = append(combined.llmreq, other.llmreq...)
	combined.a2areqs = append(combined.a2areqs, other.a2areqs...)
	// Sum other values
	combined.llmreqtokens += other.llmreqtokens
	combined.llmresptokens += other.llmresptokens
	combined.llmreqbytes += other.llmreqbytes
	combined.llmrespbytes += other.llmrespbytes
	combined.a2areqtokens += other.a2areqtokens
	combined.a2aresptokens += other.a2aresptokens
	combined.a2areqbytes += other.a2areqbytes
	combined.a2arespbytes += other.a2arespbytes
	return combined
}

// Tokens counts the number of tokens in a given text using the tiktoken library.
func Tokens(text string) (int, error) {
	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())
	te, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return 0, fmt.Errorf("get encoding: %v", err)
	}
	return len(te.Encode(text, nil, nil)), nil
}

// Entry is a struct that represents a single entry in the statistics.
type Entry struct {
	// Title is the name of the statistic.
	Title string

	// Value is the value of the statistic, formatted as a string.
	Value string
}

// Entries generates a list of statistics entries.
func (s *Stats) Entries() []Entry {
	return []Entry{
		{"Total LLM messages asked", fmt.Sprintf("%d", len(s.TotalLLMRequests()))},
		{"Total LLM request duration", s.TotalLLMReqDuration().String()},
		{"Total LLM tokens", fmt.Sprintf("%d", s.TotalLLMTokens())},
		{"Total LLM request tokens", fmt.Sprintf("%d", s.TotalLLMReqTokens())},
		{"Total LLM response tokens", fmt.Sprintf("%d", s.TotalLLMRespTokens())},
		{"Total LLM bytes", fmt.Sprintf("%d", s.TotalLLMBytes())},
		{"Total LLM request bytes", fmt.Sprintf("%d", s.TotalLLMReqBytes())},
		{"Total LLM response bytes", fmt.Sprintf("%d", s.TotalLLMRespBytes())},
		{"Average LLM request duration", s.AverageLLMReqDuration().String()},
		{"Average LLM request tokens", fmt.Sprintf("%.4f", s.AverageLLMReqTokens())},
		{"Average LLM response tokens", fmt.Sprintf("%.4f", s.AverageLLMRespTokens())},
		{"Average LLM request bytes", fmt.Sprintf("%.4f", s.AverageLLMReqBytes())},
		{"Average LLM response bytes", fmt.Sprintf("%.4f", s.AverageLLMRespBytes())},
		{"Total A2A messages asked", fmt.Sprintf("%d", s.A2AMessages())},
		{"Total A2A request duration", s.TotalA2AReqDuration().String()},
		{"Total A2A tokens", fmt.Sprintf("%d", s.TotalA2ATokens())},
		{"Total A2A request tokens", fmt.Sprintf("%d", s.TotalA2AReqTokens())},
		{"Total A2A response tokens", fmt.Sprintf("%d", s.TotalA2ARespTokens())},
		{"Total A2A bytes", fmt.Sprintf("%d", s.TotalA2ABytes())},
		{"Total A2A request bytes", fmt.Sprintf("%d", s.TotalA2AReqBytes())},
		{"Total A2A response bytes", fmt.Sprintf("%d", s.TotalA2ARespBytes())},
		{"Average A2A request duration", s.AverageA2AReqDuration().String()},
		{"Average A2A request tokens", fmt.Sprintf("%.4f", s.AverageA2AReqTokens())},
		{"Average A2A response tokens", fmt.Sprintf("%.4f", s.AverageA2ARespTokens())},
		{"Average A2A request bytes", fmt.Sprintf("%.4f", s.AverageA2AReqBytes())},
		{"Average A2A response bytes", fmt.Sprintf("%.4f", s.AverageA2ARespBytes())},
	}
}
