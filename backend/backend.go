package backend

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type BackendServer struct {
	url *url.URL
	AliveStatus bool
	mux sync.RWMutex
}

func NewBackendServer(url *url.URL, alive bool) *BackendServer {
	return &BackendServer{
		url: url,
		AliveStatus: alive,
	}
}

// multiple go routines can access this concurrently
func (bs *BackendServer) isAlive() bool {
	var status bool
	bs.mux.Lock()
	status = bs.AliveStatus
	bs.mux.Unlock()
	return status
}



/* Load balancer logic */

type LoadBalancer struct {
	servers []*BackendServer
	roundRobinCtr int32
}

func NewLoadBalancer(servers []*BackendServer) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		roundRobinCtr: 0,
	}
}

// increment by 1 (bounded in 0 to len-1)
func (lb *LoadBalancer) next() int {
	return int(atomic.AddInt32(&lb.roundRobinCtr, 1)) % int(len(lb.servers))
}

func (lb *LoadBalancer) nextHealthyServer() *BackendServer {
	currId := int(lb.roundRobinCtr)
	nextId := lb.next()

	// we loop once from nextId till nextId-1 (circularly)
	for !lb.servers[nextId].isAlive() && nextId != currId {
		nextId = lb.next()
	}

	if nextId != currId {
		return lb.servers[nextId]
	} else if lb.servers[currId].isAlive() {
		// additionally checking for thr current server
		return lb.servers[currId]
	}

	// no healthy server available
	return nil
}

// Polling for health status 


// entry point of the request

// writing the ServeHTTP(http.ResponseWriter, *http.Request) to implement the http.Handler Interface 
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.nextHealthyServer()

	if server == nil {
		log.Fatal("no server available")
	}

	rp := httputil.NewSingleHostReverseProxy(server.url) 
	rp.ServeHTTP(w,r)
}