package store

import (
	"time"

	"github.com/mwf/golidays/model"
)

// Store is an interface for storing an getting holidays
type Store interface {
	Set(holidays model.Holidays) error
	Get(date time.Time) (model.Holiday, bool, error)
	GetRange(from, to time.Time) (model.Holidays, error)
}
