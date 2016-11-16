package service

import (
	"time"

	"github.com/mwf/golidays/model"
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
}

// service is a simple Service interface implementation
type service struct {
	updater *Updater
	storage store.Store
	log     logger.Logger
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

	return s, nil
}

func (s *service) Run() error {
	if s.updater != nil {
		s.updater.Run()
	}
	return nil
}

func (s *service) Stop() {
	if s.updater != nil {
		s.updater.Stop()
	}
}

func (s *service) Get(date time.Time) (model.Holiday, bool, error) {
	return s.storage.Get(date)
}

func (s *service) GetRange(from, to time.Time) (model.Holidays, error) {
	return s.storage.GetRange(from, to)
}
