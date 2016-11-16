package service

import (
	"fmt"
	"time"

	"github.com/mwf/golidays/model"
	"github.com/mwf/golidays/service/backuper"
	"github.com/mwf/golidays/service/logger"
	"github.com/mwf/golidays/service/store"
)

// Service is an interface for holidays storage with optional maintenance
// (periodic updates, backups, etc.)
type Service interface {
	// Run starts periodic jobs
	Run() error
	// Stop stops all periodic jobs
	Stop()
	// Getters from Store interface
	store.HolidayGetter

	// RestoreStorage wipes storage and restores it from the last backup
	RestoreStorage() error
}

// service is a simple Service interface implementation
type service struct {
	updater  *Updater
	backuper *backuper.Backuper
	storage  store.Store
	log      logger.Logger
}

func NewService(config *Config) (Service, error) {
	config.Defaultize()
	if err := config.Validate(); err != nil {
		return nil, err
	}

	s := &service{
		storage: config.Storage,
		log:     config.Logger,
	}

	if !config.Updater.Disabled {
		updater, err := NewUpdater(config.Storage, config.Updater.Crawler, config.Updater.Period, s.log)
		if err != nil {
			return nil, err
		}
		s.updater = updater
	}

	if !config.Backuper.Disabled {
		b, err := backuper.New(
			config.Storage, config.Backuper.Period, config.Backuper.BasePath,
			config.Backuper.MaxBackups, s.log)
		if err != nil {
			return nil, err
		}
		s.backuper = b
	}

	return s, nil
}

func (s *service) Run() error {
	if s.updater != nil {
		s.updater.Run()
	}
	if s.backuper != nil {
		s.backuper.Run()
	}
	return nil
}

func (s *service) Stop() {
	if s.updater != nil {
		s.updater.Stop()
	}
	if s.backuper != nil {
		s.backuper.Stop()
	}
}

func (s *service) Get(date time.Time) (model.Holiday, bool, error) {
	return s.storage.Get(date)
}

func (s *service) GetRange(from, to time.Time) (model.Holidays, error) {
	return s.storage.GetRange(from, to)
}

func (s *service) RestoreStorage() error {
	if s.backuper == nil {
		return fmt.Errorf("backuper is disabled")
	}

	return s.backuper.RestoreStorage()
}
