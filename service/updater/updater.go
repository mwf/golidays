package updater

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
	minUpdatePeriod = time.Minute
)

// Updater performs periodic holiday updates in storage for current year
type Updater struct {
	storage store.Store
	crawler crawler.Crawler
	period  time.Duration
	logger  logger.Logger

	runOnce sync.Once
	done    chan struct{}
}

// New returns new updater instance
func New(storage store.Store, crawler crawler.Crawler, period time.Duration, log logger.Logger) (*Updater, error) {
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
	u.logger.Infof("%s started", u)
	defer u.logger.Infof("%s stopped", u)

	// perform initial update on start
	u.perform()
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

	year := startedAt.Year()
	if startedAt.Month() >= time.November {
		// start scraping next year in november
		year += 1
	}
	h, err := u.crawler.ScrapeYear(year)
	if err != nil {
		u.logger.Errorf("crawler.ScrapeYear error: %s", err)
		return
	}

	if err := u.storage.Set(h); err != nil {
		u.logger.Errorf("storage.Set error: %s", err)
		return
	}
}
