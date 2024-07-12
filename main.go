package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/ChiragRajput101/loadbalancer/algorithm"
	"github.com/ChiragRajput101/loadbalancer/backend"
	"github.com/ChiragRajput101/loadbalancer/config"
	"github.com/ChiragRajput101/loadbalancer/health"
)
 

type LoadBalancer struct {
	// for loading in the required configuration from the config.yaml file
	config *config.Config
	// ServerList -> for final integration of servers[], algorithm, Name, HealthChecker
	// mapping the matcher to this final configuration 
	ServerMap map[string]*config.ServerList
}

func NewLoadBalancer(cfg *config.Config) *LoadBalancer {
	serverMap := make(map[string]*config.ServerList)

	// for each service defined in the config.yaml file
	for _,service := range cfg.Services {

		// setting up the replicas 
		servers := make([]*backend.Server,0)

		// for each replica in a service
		for _, replica := range service.Replicas {
			u, err := url.Parse(replica.URL)
			if err != nil {
				log.Fatal(err)	
			}
			revProxy := httputil.NewSingleHostReverseProxy(u)

			// store in the replica, initialised as the server with URL and the httputil.ReverseProxy
			servers = append(servers, &backend.Server{
				URL: u,
				RevProxy: revProxy,
			})
		}
		
		// init the healthChecker for all servers(replicas)
		healthChecker := health.NewHealthChecker(servers)	

		// mapping the matcher to the final configuration of a service
		serverMap[service.Matcher] = &config.ServerList{
			Servers: servers,
			Name: service.Name,
			Strategy: algorithm.LoadStrategy(service.Strategy),
			Hc: healthChecker,	
		}
	}

	// starting the health check for all provided matchers
	for _,s := range serverMap {
		go s.Hc.Start()
	}

	return &LoadBalancer{
		config: cfg,
		ServerMap: serverMap,
	}
}

func (lb *LoadBalancer) findServiceList(path string) (*config.ServerList, error) {
	for matcher, s := range lb.ServerMap {
		if strings.HasPrefix(path, matcher) {
			fmt.Printf("Found service '%s' matching the request\n", s.Name)
			return s, nil
		}
	}
	return nil, fmt.Errorf("could not find a matcher for url: '%s'", path)
}

// implementing the http Handler interface so that LoadBalancer can be passed as the Handler to http Server
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received new request: url='%s' \n", r.Host)
	sl, err := lb.findServiceList(r.URL.Path)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	next, err := sl.Strategy.Next(sl.Servers)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Forwarding to the server='%s' \n", next.URL.Host)
	next.ForwardRequest(w, r)
}

func main() {
	file,err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	cfg,err := config.LoadConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	lb := NewLoadBalancer(cfg)

	server := http.Server{
		Addr: ":5000",
		Handler: lb,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}