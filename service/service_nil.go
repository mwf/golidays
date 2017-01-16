package service

import (
	"time"

	"github.com/mwf/golidays/model"
)

// nilService is a Service doing nothing
type nilService struct{}

// NewNilService returns Service instance, doing nothing
func NewNilService() Service {
	return &nilService{}
}

func (s *nilService) Run() error {
	return nil
}

func (s *nilService) Stop() {}

func (s *nilService) Get(date time.Time) (model.Holiday, bool, error) {
	return model.Holiday{}, false, nil
}

func (s *nilService) GetRange(from, to time.Time) (model.Holidays, error) {
	return nil, nil
}

func (s *nilService) RestoreStorage() error {
	return nil
}
