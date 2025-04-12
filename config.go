package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const defaultListenAddr = ":9999"
const defaultTorAddr = "localhost:9050"
const defaultTimeout = time.Second * 10
const defaultCheckInterval = time.Second * 30
const defaultLogLevel = "error"
const defaultLogFormat = "fmt"
const targetTypeHTTP = "http"
const targetTypeTCP = "tcp"

type Config struct {
	ListenAddr         string        `yaml:"listen_addr,omitempty"`
	TorAddr            string        `yaml:"tor_addr,omitempty"`
	Timeout            time.Duration `yaml:"timeout,omitempty"`
	CheckInterval      time.Duration `yaml:"check_interval,omitempty"`
	LogLevel           string        `yaml:"log_level"`
	LogFormat          string        `yaml:"log_format"`
	Targets            []Target      `yaml:"targets,omitempty"`
	InsecureSkipVerify bool          `yaml:"insecure_skip_ssl_verify"`
}

type Target struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	URL  string `yaml:"url"`
}

func NewConfig() *Config {
	return &Config{
		ListenAddr:         defaultListenAddr,
		TorAddr:            defaultTorAddr,
		Timeout:            defaultTimeout,
		CheckInterval:      defaultCheckInterval,
		LogLevel:           defaultLogLevel,
		LogFormat:          defaultLogFormat,
		InsecureSkipVerify: false,
	}
}

func LoadConfig(path *string) (*Config, error) {
	cfg := NewConfig()

	bytes, err := os.ReadFile(*path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
