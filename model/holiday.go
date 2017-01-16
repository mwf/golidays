package model

import (
	"time"
)

// HolidayType describes different holiday types, e.g. weekends, preholiday days.
type HolidayType string

const (
	TypeWeekend    = HolidayType("weekend")    // just an ordinary weekend day
	TypeHoliday    = HolidayType("holiday")    // a holiday, real or shifted from the weekend
	TypePreholiday = HolidayType("preholiday") // a preholiday day considered contracted
)

type Holiday struct {
	Date time.Time   `json:"date" yaml:"date"`
	Type HolidayType `json:"type" yaml:"type"`
}

type Holidays []Holiday

type HolidaysByDate Holidays

func (h HolidaysByDate) Len() int           { return len(h) }
func (h HolidaysByDate) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h HolidaysByDate) Less(i, j int) bool { return h[i].Date.Before(h[j].Date) }

// NewDay returns time.Time instance, representing day
func NewDay(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
