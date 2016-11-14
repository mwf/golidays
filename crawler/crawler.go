package crawler

import (
	"github.com/mwf/golidays/model"
)

// Crawler is an interface for parsing holidays from different websites
type Crawler interface {
	ScrapeYear() (model.Holidays, error)
}
