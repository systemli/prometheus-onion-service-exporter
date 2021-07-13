package main

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

const defaultListenAddr = ":9999"
const defaultTorAddr = "localhost:9050"
const defaultTimeout = time.Second * 10
const defaultCheckInterval = time.Second * 30

type Config struct {
	ListenAddr    string        `yaml:"listen_addr,omitempty"`
	TorAddr       string        `yaml:"tor_addr,omitempty"`
	Timeout       time.Duration `yaml:"timeout,omitempty"`
	CheckInterval time.Duration `yaml:"check_interval,omitempty"`
	Targets       []Target      `yaml:"targets,omitempty"`
}

type Target struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func NewConfig() *Config {
	return &Config{
		ListenAddr:    defaultListenAddr,
		TorAddr:       defaultTorAddr,
		Timeout:       defaultTimeout,
		CheckInterval: defaultCheckInterval,
	}
}

func LoadConfig(path *string) (error, *Config) {
	cfg := NewConfig()

	bytes, err := ioutil.ReadFile(*path)
	if err != nil {
		return err, cfg
	}

	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return err, cfg
	}

	return nil, cfg
}
