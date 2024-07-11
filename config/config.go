package config

import (
	"io"

	"github.com/ChiragRajput101/loadbalancer/algorithm"
	"github.com/ChiragRajput101/loadbalancer/backend"
	"github.com/ChiragRajput101/loadbalancer/health"
	"gopkg.in/yaml.v3"
)

// reads from the config file
type Config struct {
	Services []backend.Service `yaml:"services"`
	Strategy string `yaml:"strategy"`
}

type ServerList struct {
	// Servers are the replicas
	Servers []*backend.Server

	// Name of the service
	Name string

	// Strategy defines how the server list is load balanced.
	// It can never be 'nil', it should always default to a 'RoundRobin' version.
	Strategy algorithm.Algorithm

	// Health checker for the servers
	Hc *health.HealthChecker
}

func LoadConfig(r io.Reader) (*Config, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	conf := Config{}
	if err := yaml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}