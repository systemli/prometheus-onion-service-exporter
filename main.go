package main

import (
	"crypto/tls"
	"flag"
	"net/http"
	"net/url"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

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
		log.WithError(err).Fatal("unable to load the config file")
	}

	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.WithError(err).Error("unable to parse the log level")
	} else {
		log.SetLevel(lvl)
	}

	if cfg.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
}

func main() {
	log.Infof("started Onion Service Exporter on %s", cfg.ListenAddr)

	go func() {
		for {
			log.Debug("start loop to check the targets")

			func() {
				var wg sync.WaitGroup

				for _, target := range cfg.Targets {
					log.Debugf(`check target "%s"`, target.Name)

					wg.Add(1)

					switch target.Type {
					case targetTypeHTTP:
						checkHTTP(target, &wg)
					case targetTypeTCP:
						checkTCP(target, &wg)
					default:
						log.Errorf(`unsupported scheme "%s" for target "%s"`, target.Type, target.Name)
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
		log.WithError(err).Error("unable to parse the url")
		return
	}

	dialer, err := proxy.SOCKS5("tcp", cfg.TorAddr, nil, proxy.Direct)
	if err != nil {
		log.WithError(err).Error("failed to initialize the proxy")
	}

	transport := &http.Transport{Dial: dialer.Dial}
	if cfg.IgnoreSSL {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{Transport: transport, Timeout: cfg.Timeout}
	up := 0.0
	start := time.Now()

	res, err := client.Get(uri.String())
	if err != nil {
		log.WithError(err).WithField("url", uri.String()).Warn("unable to get the url")
	} else {
		if res.StatusCode == http.StatusOK {
			up = 1.0
		}
		defer res.Body.Close()
	}

	statuses[target.Name] = OnionStatus{
		Up:      up,
		Host:    uri.Host,
		Type:    targetTypeHTTP,
		Latency: time.Since(start).Seconds(),
	}

	log.Debugf(`finished check for "%s"`, target.Name)
}

func checkTCP(target Target, wg *sync.WaitGroup) {
	wg.Done()
	uri, err := url.Parse(target.URL)
	if err != nil {
		log.WithError(err).Error("unable to parse the url")
		return
	}

	dialer, err := proxy.SOCKS5("tcp", cfg.TorAddr, nil, proxy.Direct)
	if err != nil {
		log.WithError(err).Error("failed to initialize the proxy")
	}

	up := 0.0
	start := time.Now()
	conn, err := dialer.Dial("tcp", uri.Host)
	if err != nil {
		log.WithError(err).WithField("url", uri.String()).Warn("unable to get the url")
	} else {
		up = 1.0
		defer conn.Close()
	}

	statuses[target.Name] = OnionStatus{
		Up:      up,
		Host:    uri.Host,
		Type:    targetTypeTCP,
		Latency: time.Since(start).Seconds(),
	}
}
