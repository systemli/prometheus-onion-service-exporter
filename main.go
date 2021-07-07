package main

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/proxy"
)

var (
	client *http.Client
)

func main() {
	err, cfg := LoadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	dialer, err := proxy.SOCKS5("tcp", cfg.TorAddr, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{Dial: dialer.Dial}
	client = &http.Client{Transport: transport, Timeout: cfg.Timeout}
	ticker := time.NewTicker(cfg.CheckInterval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				func() {
					var wg sync.WaitGroup

					for _, target := range cfg.Targets {
						wg.Add(1)
						go checkOnionService(target, &wg)
					}
					wg.Wait()
				}()
			}
		}
	}()

	prometheus.MustRegister(NewOnionCollector())
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Started Onion Service Exporter on %s", cfg.ListenAddr)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func checkOnionService(target Target, wg *sync.WaitGroup) {
	defer wg.Done()

	uri, err := url.Parse(target.URL)
	if err != nil {
		log.Println(err)
		return
	}

	up := 0.0
	start := time.Now()
	res, err := client.Get(uri.String())
	if err != nil {
		log.Println(err)
	} else {
		if res.StatusCode == http.StatusOK {
			up = 1.0
		}
	}

	statuses[target.Name] = OnionStatus{
		Up:      up,
		Host:    uri.Host,
		Latency: time.Since(start).Seconds(),
	}
}
