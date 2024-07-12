package algorithm

import (
	"fmt"
	"sync"

	"github.com/ChiragRajput101/loadbalancer/backend"
)

type Algorithm interface {
	// gets the next healthy server, else nil
	Next([]*backend.Server) (*backend.Server, error)
}

var strategy map[string] func() Algorithm


type RoundRobin struct {
	mux sync.RWMutex
	ctr int
}

func (r *RoundRobin) Next(servers []*backend.Server) (*backend.Server, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	curr := r.ctr
	r.ctr = ((r.ctr + 1) % int(len(servers)))

	for !servers[r.ctr].IsAlive() && r.ctr != curr {
		r.ctr = ((r.ctr + 1) % int(len(servers)))
	}

	if r.ctr != curr {
		return servers[r.ctr], nil
	} else if servers[curr].IsAlive() {
		return servers[curr], nil
	}

	return nil, fmt.Errorf("no healthy server available")
}

func LoadStrategy(algo string) Algorithm {
	strategy = make(map[string]func() Algorithm)

	// init algo
	strategy["RoundRobin"] = func() Algorithm {
		return &RoundRobin{
			ctr: 0,
		}
	}

	st, ok := strategy[algo]
	if !ok {
		fmt.Println("no such strategy, falling back to Round Robin")
		return strategy["RoundRobin"]()
	}
	return st()
}
