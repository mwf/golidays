package service

import (
	"github.com/mwf/golidays/service/store"
)

// Service is an interface for holidays storage with optional maintenance
// (periodic updates, backups, etc.)
type Service interface {
	// Run starts periodic jobs
	Run() error
	// Getters from Store interface
	store.HolidayGetter
}
