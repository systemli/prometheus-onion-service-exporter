package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/proxy"
)

var cfg *Config

func init() {
	configFile := flag.String("c", "config.yml", "Path to your config file")
	flag.Parse()

	var err error
	err, cfg = LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Printf("Started Onion Service Exporter on %s", cfg.ListenAddr)

	go func() {
		for {
			func() {
				var wg sync.WaitGroup

				for _, target := range cfg.Targets {
					wg.Add(1)

					switch target.Type {
					case targetTypeHTTP:
						checkHTTP(target, &wg)
					case targetTypeTCP:
						checkTCP(target, &wg)
					default:
						log.Printf(`Unsupported scheme "%s"\n`, target.Type)
					}
				}
				wg.Wait()
			}()
			<-time.After(cfg.CheckInterval)
		}
	}()

	prometheus.MustRegister(NewOnionCollector())
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func checkHTTP(target Target, wg *sync.WaitGroup) {
	wg.Done()
	uri, err := url.Parse(target.URL)
	if err != nil {
		log.Println(err)
		return
	}

	dialer, err := proxy.SOCKS5("tcp", cfg.TorAddr, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{Dial: dialer.Dial}
	client := &http.Client{Transport: transport, Timeout: cfg.Timeout}
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
		Type:    targetTypeHTTP,
		Latency: time.Since(start).Seconds(),
	}
}

func checkTCP(target Target, wg *sync.WaitGroup) {
	wg.Done()
	uri, err := url.Parse(target.URL)
	if err != nil {
		log.Println(err)
		return
	}

	dialer, err := proxy.SOCKS5("tcp", cfg.TorAddr, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	up := 0.0
	start := time.Now()
	_, err = dialer.Dial("tcp", uri.Host)
	if err != nil {
		log.Println(err)
	} else {
		up = 1.0
	}

	statuses[target.Name] = OnionStatus{
		Up:      up,
		Host:    uri.Host,
		Type:    targetTypeTCP,
		Latency: time.Since(start).Seconds(),
	}
}
