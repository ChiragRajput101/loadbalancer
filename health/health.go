package health

import (
	"fmt"
	"net"
	"time"

	"github.com/ChiragRajput101/loadbalancer/backend"
)

type HealthChecker struct {
	Servers []*backend.Server
}

func NewHealthChecker(servers []*backend.Server) *HealthChecker {
	return &HealthChecker{
		Servers: servers,
	}
}

// running a seperate go routine for health check
func (h *HealthChecker) Start() {

	fmt.Printf("starting up the health check")
	// Ticker containing a channel that will send the current time on the channel after each tick
	// the period of the ticks is specified by the duration argument.
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case _ = <-ticker.C:
			for _,server := range h.Servers {
				go check(server)
			}
		}
	}
}

func check(server *backend.Server) {
	_, err := net.DialTimeout("tcp", server.URL.Host, time.Second * 3)
	if err != nil {
		fmt.Printf("could not connect to server %v", server.URL.Host)
		server.SetAliveStatus(false)
	} else {
		server.SetAliveStatus(true)
	}
}	
