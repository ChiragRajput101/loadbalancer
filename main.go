package main

import (
	"log"
	"net/http"

	"github.com/ChiragRajput101/loadbalancer/backend"
	"github.com/ChiragRajput101/loadbalancer/util"
)

func main() {
	
	serverPool := util.SpinUpServers("5000",5)

	lb := backend.NewLoadBalancer(serverPool)
	// lb implements the interface http.Handler
	http.Handle("/", lb)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("base server issue")
	}
}