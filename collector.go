package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type OnionStatus struct {
	Up      float64
	Host    string
	Latency float64
}

var statuses map[string]OnionStatus

type OnionCollector struct {
	Up      *prometheus.Desc
	Latency *prometheus.Desc
}

func NewOnionCollector() *OnionCollector {
	statuses = make(map[string]OnionStatus)

	return &OnionCollector{
		Up:      prometheus.NewDesc("onion_service_up", "", []string{"name", "address"}, nil),
		Latency: prometheus.NewDesc("onion_service_latency", "", []string{"name", "address"}, nil),
	}
}

func (oc *OnionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- oc.Up
	ch <- oc.Latency
}

func (oc *OnionCollector) Collect(ch chan<- prometheus.Metric) {
	for name, status := range statuses {
		ch <- prometheus.MustNewConstMetric(
			oc.Up,
			prometheus.GaugeValue,
			status.Up,
			name,
			status.Host,
		)
		if status.Up != 0 {
			ch <- prometheus.MustNewConstMetric(
				oc.Latency,
				prometheus.GaugeValue,
				status.Latency,
				name,
				status.Host,
			)
		}
	}
}
