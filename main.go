package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var MongoUser string
var MongoPass string
var MongoNodes []string

func main() {
	MongoUser = os.Getenv("MONGO_USER")
	MongoPass = os.Getenv("MONGO_PASSWORD")

	nodes := os.Getenv("MONGO_NODES")
	if nodes == "" {
		panic("no nodes to probe")
	}

	MongoNodes = strings.Split(nodes, ",")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8100"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	lvl, err := logrus.ParseLevel(logLevel)
	logrus.SetLevel(lvl)
	logrus.Info("Log level set to " + logLevel)

	// Disable go_exporter default metrics
	r := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = r
	prometheus.DefaultGatherer = r
	initMetrics()

	check()
	go tick()

	http.Handle("/metrics", promhttp.Handler())
	logrus.Info("Starting metrics server on port " + port)
	err = http.ListenAndServe(":" + port, nil)
	if err != nil {
		logrus.Error(err)
		return
	}
}