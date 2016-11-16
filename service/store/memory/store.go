package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/mwf/golidays/model"
	"github.com/mwf/golidays/service/store"
)

// Store is a simple in-memory storage
type Store struct {
	byDate map[time.Time]*model.Holiday
	mu     sync.RWMutex
}

// check if Store implements Store interface
var _ store.Store = New()

func New() *Store {
	return &Store{
		byDate: make(map[time.Time]*model.Holiday),
	}
}

// Set set's holidays to in-memory storage
func (s *Store) Set(holidays model.Holidays) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, holiday := range holidays {
		s.byDate[holiday.Date] = &holidays[i]
	}

	return nil
}

// Get finds holiday by date and returns it.
// If nothing found - returned bool value is false
func (s *Store) Get(date time.Time) (model.Holiday, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if holiday, ok := s.byDate[date]; ok {
		return *holiday, ok, nil
	}

	return model.Holiday{}, false, nil
}

// GetRange returns holidays between 'from' and 'to' dates
func (s *Store) GetRange(from, to time.Time) (model.Holidays, error) {
	return s.getRangeNaive(from, to)
}

// Naive implementation - simply range over all dates in range
// Returns empty slice if no holidays found
func (s *Store) getRangeNaive(from, to time.Time) (model.Holidays, error) {
	if to.Before(from) {
		return nil, fmt.Errorf("invalid range: %s > %s", from, to)
	}

	holidays := model.Holidays{}
	s.mu.RLock()
	defer s.mu.RUnlock()

	for from.Before(to) {
		if h, ok := s.byDate[from]; ok {
			holidays = append(holidays, *h)
		}
		from = from.Add(24 * time.Hour)
	}
	// check the last day
	if h, ok := s.byDate[to]; ok {
		holidays = append(holidays, *h)
	}

	return holidays, nil
}

// Dump returns all holidays from store in random order
func (s *Store) Dump() model.Holidays {
	s.mu.RLock()
	defer s.mu.RUnlock()

	holidays := make(model.Holidays, 0, len(s.byDate))
	for _, h := range s.byDate {
		holidays = append(holidays, *h)
	}

	return holidays
}
