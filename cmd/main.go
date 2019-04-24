package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/mwf/golidays/crawler"
	"github.com/mwf/golidays/service"
	"github.com/mwf/golidays/service/store/memory"
	"github.com/sirupsen/logrus"
)

func waitInterrupt() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case s := <-c:
		logrus.Infof("Recieved %s, exiting...", s)
		os.Exit(0)
	}
}

func main() {
	c := crawler.NewConsultantRu()
	storage := memory.New()
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	config := &service.Config{
		Updater: service.UpdaterConfig{
			Crawler: c,
			Period:  5 * time.Minute,
		},
		Backuper: service.BackuperConfig{
			Period:   1 * time.Minute,
			BasePath: "./var",
		},
		Storage: storage,
		Logger:  logger,
	}

	srv, err := service.New(config)
	if err != nil {
		logger.Warnf("error initializing service: %s", err)
		os.Exit(1)
	}
	if err := srv.RestoreStorage(); err != nil {
		logger.Warnf("error restoring storage: %s", err)
	}
	srv.Run()

	waitInterrupt()
}
