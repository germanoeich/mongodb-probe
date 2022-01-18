package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ErrorCounter prometheus.Counter
	StatusGauge *prometheus.GaugeVec
	ReplicationLagGauge *prometheus.GaugeVec
	UpGauge prometheus.Gauge
)

func initMetrics() {
	ErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mongodb_probe_error",
		Help: "The total number of errors when processing requests",
	})

	StatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mongodb_probe_status",
		Help:        "Gauge with node state. https://docs.mongodb.com/manual/reference/replica-states/",
	}, []string{"node"})

	ReplicationLagGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mongodb_probe_replication_lag",
		Help:        "Gauge with node replication lag in ms",
	}, []string{"node"})

	UpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "mongodb_probe_up",
		Help:        "Gauge with uptime state - always 1 if application is alive",
	})
	UpGauge.Set(1)
}