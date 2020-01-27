package main

import (
	"context"
	"flag"
	"github.com/ducc/profile-collector/collector"
	"github.com/ducc/profile-collector/collector/config"
	"github.com/ducc/profile-collector/collector/scraper"
	"github.com/ducc/profile-collector/metrics"
	"github.com/ducc/profile-collector/store/client"
	"github.com/sirupsen/logrus"
)

var (
	logLevel   string
	configPath string
)

func init() {
	flag.StringVar(&logLevel, "level", "debug", "logrus logging level")
	flag.StringVar(&configPath, "config", "", "config path") // todo
}

func main() {
	logging()
	go metrics.Start()

	ctx := context.Background()

	conf, err := config.New(configPath)
	if err != nil {
		logrus.WithError(err).Fatal("getting config")
	}

	scrap, err := scraper.NewHTTP()
	if err != nil {
		logrus.WithError(err).Fatal("creating scraper")
	}

	store, err := client.New(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("creating store client")
	}

	if err := collector.Run(ctx, conf, scrap, store); err != nil {
		logrus.WithError(err).Fatal("running collector")
	}
}

func logging() {
	if ll, err := logrus.ParseLevel(logLevel); err != nil {
		logrus.WithError(err).Fatal("unable to parse logrus logging level")
	} else {
		logrus.SetLevel(ll)
	}
}
