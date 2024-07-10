package util

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/ChiragRajput101/loadbalancer/backend"
)

func SpinUpServers(start string, count int) []*backend.BackendServer {
	list := []*url.URL{}

	port,_ := strconv.Atoi(start)

	for i:=0;i<count;i++ {
		s := fmt.Sprintf("http://localhost:%v",port)
		u,e := url.Parse(s)
		if e != nil {
			log.Fatal("not able to parse URL")
		}
		list = append(list, u)
		port++
	}

	servers := []*backend.BackendServer{}

	for _,v := range list {
		s := backend.NewBackendServer(v,true)
		servers = append(servers, s)
	}

	return servers
}