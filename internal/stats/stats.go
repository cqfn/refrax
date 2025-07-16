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
	// a2areqs []time.Duration

	// a2areqtokens is an int that stores the number of tokens used in A2A requests.
	// a2areqtokens int

	// a2areqbytes is an int that stores the number of bytes used in A2A requests.
	// a2areqbytes int

	// a2aresptokens is an int that stores the number of tokens received in A2A responses.
	// a2aresptokens int

	// a2arespbytes is an int that stores the number of bytes received in A2A responses.
	// a2arespbytes int
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
// @toto #60:90min Implement A2AReq Method.
// This method is currently a placeholder and should be implemented to record
// A2A request statistics.
func (s *Stats) A2AReq(_ time.Duration, _, _, _, _ int) {
	panic("not implemented")
}

// LLMRequests retrieved all request timings.
// @todo #60:90min Print Entire Statistics.
// Currently, this function only returns the timings of LLM requests.
// It should be extended to return all statistics, including tokens and bytes.
// Don't forget to update std and csv writers.
func (s *Stats) LLMRequests() []time.Duration {
	duplicate := make([]time.Duration, len(s.llmreq))
	copy(duplicate, s.llmreq)
	return duplicate
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
