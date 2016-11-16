package service

import (
	"fmt"
	"time"

	"github.com/mwf/golidays/crawler"
	"github.com/mwf/golidays/service/logger"
	"github.com/mwf/golidays/service/store"
	"github.com/mwf/golidays/service/store/memory"
)

const (
	defaultUpdatePeriod = 24 * time.Hour
	defaultBackupPeriod = 7 * 24 * time.Hour
	defaulMaxBackups    = 4
)

// Config is a service configuration data struct
type Config struct {
	Updater  UpdaterConfig
	Backuper BackuperConfig
	Storage  store.Store
	Logger   logger.Logger
}

type UpdaterConfig struct {
	Disabled bool
	Period   time.Duration
	Crawler  crawler.Crawler
}

type BackuperConfig struct {
	Disabled   bool
	BasePath   string
	Period     time.Duration
	MaxBackups int
}

// Defaultize sets default values for some config values
func (c *Config) Defaultize() {
	if c.Updater.Period == 0 {
		c.Updater.Period = defaultUpdatePeriod
	}

	if c.Backuper.Period == 0 {
		c.Backuper.Period = defaultBackupPeriod
	}
	if c.Backuper.MaxBackups == 0 {
		c.Backuper.MaxBackups = defaulMaxBackups
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
		return fmt.Errorf("config.Updater.Crawler is nil")
	}

	if !c.Backuper.Disabled && c.Backuper.BasePath == "" {
		return fmt.Errorf("config.Backuper.BasePath is empty")
	}

	return nil
}
