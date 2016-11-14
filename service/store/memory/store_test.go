package memory

import (
	"reflect"
	"testing"
	"time"

	"github.com/mwf/golidays/model"
)

func newHolidays() model.Holidays {
	holidays := model.Holidays{}
	for i := 25; i <= 31; i++ {
		h := model.Holiday{
			Date: time.Date(1977, 5, i, 0, 0, 0, 0, time.UTC),
			Type: model.TypeHoliday,
		}
		holidays = append(holidays, h)
	}

	return holidays
}

func TestSet(t *testing.T) {
	store := New()
	holidays := newHolidays()

	err := store.Set(holidays)
	if err != nil {
		t.Fatalf("Set failed: %s", err)
	}

	for _, h := range holidays {
		storedH, ok := store.byDate[h.Date]
		if !ok {
			t.Errorf("holiday %#v not found", h)
		}
		if *storedH != h {
			t.Errorf("stored data %v != original %v", storedH, h)
		}
	}
}

func TestGet(t *testing.T) {
	store := New()
	holidays := newHolidays()

	err := store.Set(holidays)
	if err != nil {
		t.Fatalf("Set failed")
	}

	for _, h := range holidays {
		storedH, ok, err := store.Get(h.Date)
		if err != nil {
			t.Fatalf("Get failed: %s", err)
		}
		if !ok {
			t.Errorf("holiday %#v not found", h)
		}
		if storedH != h {
			t.Errorf("stored data %#v != original %#v", storedH, h)
		}
	}

	_, ok, err := store.Get(time.Now())
	if err != nil {
		t.Fatalf("Get failed: %s", err)
	}
	if ok {
		t.Errorf("Nonexisting date found")
	}
}

func TestRange_Failed(t *testing.T) {
	store := New()

	_, err := store.GetRange(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC))
	if err == nil {
		t.Fatalf("Error should not be empty")
	}
}

func TestRange_No(t *testing.T) {
	store := New()

	holidays, err := store.GetRange(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("GetRange failed: %s", err)
	}

	if len(holidays) != 0 {
		t.Fatalf("holidays should be empty: %#v", holidays)
	}
}

func TestRange_All(t *testing.T) {
	store := New()
	holidays := newHolidays()

	err := store.Set(holidays)
	if err != nil {
		t.Fatalf("Set failed")
	}

	storedH, err := store.GetRange(time.Date(1977, 5, 25, 0, 0, 0, 0, time.UTC), time.Date(1977, 5, 31, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("GetRange failed: %s", err)
	}

	if !reflect.DeepEqual(storedH, holidays) {
		t.Errorf("stored data %#v != original %#v", storedH, holidays)
	}
}
