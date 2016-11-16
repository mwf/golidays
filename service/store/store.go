package store

import (
	"time"

	"github.com/mwf/golidays/model"
)

// Store is an interface for storing an getting holidays
type Store interface {
	HolidaySetter
	HolidayGetter
	HolidayDumper
}

type HolidayGetter interface {
	Get(date time.Time) (model.Holiday, bool, error)
	GetRange(from, to time.Time) (model.Holidays, error)
}

type HolidaySetter interface {
	Set(holidays model.Holidays) error
}

type HolidayDumper interface {
	// Dump returns all items in store
	Dump() model.Holidays
}
