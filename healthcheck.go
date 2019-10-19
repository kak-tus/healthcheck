package healthcheck

import (
	"fmt"
	"net/http"
)

type Handler struct {
	http.ServeMux
}

// State type
type State int

const (
	_ = iota
	// StatePassing - check in passing state
	StatePassing State = iota
	// StateWarning - check in warning state
	StateWarning State = iota
	// StateCritical - check in critical state
	StateCritical State = iota
)

var stateMap = map[State]int{
	1: 200,
	2: 429,
	3: 500,
}

func NewHandler() *Handler {
	return &Handler{}
}

// Add add new HTTP healthcheck
func (h *Handler) Add(path string, f func() (State, string)) {
	h.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		state, text := f()

		if state != StatePassing {
			w.WriteHeader(stateMap[state])
		}

		_, _ = fmt.Fprint(w, text)
	})
}

// AddReq add new HTTP healthcheck with http.Request parameter
// to allow get some data from request
func (h *Handler) AddReq(path string, f func(*http.Request) (State, string)) {
	h.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		state, text := f(r)

		if state != StatePassing {
			w.WriteHeader(stateMap[state])
		}

		_, _ = fmt.Fprint(w, text)
	})
}
