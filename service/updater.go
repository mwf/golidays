package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/mwf/golidays/crawler"
	"github.com/mwf/golidays/service/logger"
	"github.com/mwf/golidays/service/store"
)

const (
	// no reason to keep it higher
	minUpdatePeriod     = time.Minute
	defaultUpdatePeriod = 24 * time.Hour
)

// Updater performs periodic updates of holidays in storage
type Updater struct {
	storage store.Store
	crawler crawler.Crawler
	period  time.Duration
	logger  logger.Logger

	runOnce sync.Once
	done    chan struct{}
}

// NewUpdater returns new updater instance
func NewUpdater(storage store.Store, crawler crawler.Crawler, period time.Duration, log logger.Logger) (*Updater, error) {
	if period < minUpdatePeriod {
		return nil, fmt.Errorf("period is too low: %s < %s", period, minUpdatePeriod)
	}

	return &Updater{
		storage: storage,
		crawler: crawler,
		period:  period,
		logger:  log,
		done:    make(chan struct{}),
	}, nil
}

func (u *Updater) String() string {
	return fmt.Sprintf("Updater {period: %s}", u.period)
}

// Run runs asynchronous update loop. Multiple calls do nothing - the loop started
// exactly once.
func (u *Updater) Run() {
	u.runOnce.Do(func() {
		go u.loop()
	})
}

// Stop stops update loop.
func (u *Updater) Stop() {
	u.runOnce.Do(func() {
		close(u.done)
	})
}

func (u *Updater) loop() {
	u.logger.Debugf("started %s", u)

	for {
		select {
		case <-time.After(u.period):
			u.perform()
		case <-u.done:
			return
		}
	}
}

func (u *Updater) perform() {
	startedAt := time.Now()

	u.logger.Debugf("perform %s", u)
	defer func() {
		u.logger.Infof("perform finished in %s", time.Now().Sub(startedAt))
	}()

	h, err := u.crawler.ScrapeYear(startedAt.Year())
	if err != nil {
		u.logger.Errorf("crawler.ScrapeYear error: %s", err)
		return
	}

	if err := u.storage.Set(h); err != nil {
		u.logger.Errorf("storage.Set error: %s", err)
		return
	}
}