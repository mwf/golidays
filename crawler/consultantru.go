package crawler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mwf/golidays/model"
)

const (
	yearURL = "http://www.consultant.ru/law/ref/calendar/proizvodstvennye/%d/"
)

var (
	ruMonths = map[string]int{
		"январь":   1,
		"февраль":  2,
		"март":     3,
		"апрель":   4,
		"май":      5,
		"июнь":     6,
		"июль":     7,
		"август":   8,
		"сентябрь": 9,
		"октябрь":  10,
		"ноябрь":   11,
		"декабрь":  12,
	}
)

type ConsultantRu struct{}

func NewConsultantRu() *ConsultantRu {
	return &ConsultantRu{}
}

func (c *ConsultantRu) ScrapeYear(year int) (model.Holidays, error) {
	doc, err := goquery.NewDocument(fmt.Sprintf(yearURL, year))
	if err != nil {
		return nil, err
	}

	months, err := c.getMonthTablesOrdered(doc)
	if err != nil {
		return nil, err
	}

	// In avarage there are no more than 128 holidays per year
	holidays := make(model.Holidays, 0, 128)
	for monthN, month := range months {
		monthN += 1

		var monthError error
		month.Find("td").Each(func(i int, s *goquery.Selection) {
			if !s.HasClass("weekend") && !s.HasClass("preholiday") {
				return
			}

			dayS := strings.TrimSuffix(s.Text(), "*")
			day, err := strconv.ParseInt(dayS, 10, 32)
			if err != nil {
				monthError = fmt.Errorf("can't parse day from '%s' for month %d", dayS, monthN)
				return
			}

			holiday := model.Holiday{
				Date: time.Date(year, time.Month(monthN), int(day), 0, 0, 0, 0, time.UTC),
			}
			switch {
			case s.HasClass("weekend"):
				if holiday.Date.Weekday() == time.Saturday || holiday.Date.Weekday() == time.Sunday {
					holiday.Type = model.TypeWeekend
				} else {
					holiday.Type = model.TypeHoliday
				}
			case s.HasClass("preholiday"):
				holiday.Type = model.TypePreholiday
			}
			holidays = append(holidays, holiday)
		})

		if monthError != nil {
			return nil, monthError
		}
	}

	// The date should be already sorted, but let's sort it for sure.
	sort.Sort(model.HolidaysByDate(holidays))
	return holidays, nil
}

func (c *ConsultantRu) getMonthTablesOrdered(doc *goquery.Document) (months []*goquery.Selection, err error) {
	// Find the quarter tables
	quarters := doc.Find(".calendar-table > tbody > tr:has(.quarter-title)")

	foundMonths := make([]string, 0, 12)
	quarters.Each(func(quarterN int, s *goquery.Selection) {
		if err != nil {
			return
		}

		s.Find(".month-block > table").Each(func(monthN int, s *goquery.Selection) {
			// check if every month is parsed
			if err != nil {
				return
			}

			month := strings.ToLower(s.Find("th.month").Text())
			monthNumber, ok := ruMonths[month]
			if !ok {
				err = fmt.Errorf("month '%s' does not exist", month)
				return
			}

			parsedMonthN := (3*quarterN + monthN + 1)
			if monthNumber != parsedMonthN {
				err = fmt.Errorf("month '%s' out of order - #%d, must be #%d", month, parsedMonthN, monthNumber)
				return
			}

			months = append(months, s)
			foundMonths = append(foundMonths, month)
		})
	})

	if err != nil {
		return nil, err
	}

	if len(months) != 12 {
		err = fmt.Errorf("not all months are parsed: %s", foundMonths)
		return nil, err
	}

	return
}
