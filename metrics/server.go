package metrics

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

var address string

func init() {
	flag.StringVar(&address, "metrics-address", ":8001", "address of server for prom to scrape")
}

func Start() {
	http.Handle("/metrics", promhttp.Handler())
	logrus.Debugf("starting prometheus server on %s", address)
	go func() {
		if err := http.ListenAndServe(address, nil); err != nil {
			logrus.WithError(err).Fatal("starting metrics server")
		}
	}()
}
