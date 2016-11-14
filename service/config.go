package service

import (
	"fmt"
	"time"

	"github.com/mwf/golidays/crawler"
	"github.com/mwf/golidays/service/logger"
	"github.com/mwf/golidays/service/store"
	"github.com/mwf/golidays/service/store/memory"
)

// Config is a service configuration data struct
type Config struct {
	Updater struct {
		Disabled bool
		Period   time.Duration
		Crawler  crawler.Crawler
	}
	Storage store.Store
	Logger  logger.Logger
}

// Defaultize sets default values for some config values
func (c *Config) Defaultize() {
	if c.Updater.Period == 0 {
		c.Updater.Period = defaultUpdatePeriod
	}

	if c.Storage == nil {
		c.Storage = memory.New()
	}

	if c.Logger == nil {
		c.Logger = &logger.NilLogger{}
	}
}

// Validate checks current config
func (c *Config) Validate() error {
	if !c.Updater.Disabled && c.Updater.Crawler == nil {
		return fmt.Errorf("config.Updater.Crawler can't be nil")
	}

	return nil
}
