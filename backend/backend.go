package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Replica struct {
	URL string `yaml:"url"`
}

type Service struct {
	Name string `yaml:"name"`
	Matcher string `yaml:"matcher"`
	Strategy string `yaml:"strategy"`
	Replicas []Replica `yaml:"replicas"`
}

type Config struct {
	Services []Service `yaml:"services"`
	Strategy string `yaml:"strategy"`
}

// backend server
type Server struct {
	URL *url.URL
	AliveStatus bool
	RevProxy *httputil.ReverseProxy
	mux sync.RWMutex
}

// checking the AliveStatus 
func (s *Server) IsAlive() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.AliveStatus
}

// set AliveStatus
func (s *Server) SetAliveStatus(status bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.AliveStatus = status
}

func (s *Server) ForwardRequest(w http.ResponseWriter, r *http.Request) {
	s.RevProxy.ServeHTTP(w,r)
}