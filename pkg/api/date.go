package api

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Date struct {
	t time.Time
}

func DateFromTime(t time.Time) Date {
	return Date{t}
}

func DateFromString(s string) Date {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}
	}
	return Date{t}
}

func (t *Date) Scan(src interface{}) error {
	if src == nil {
		t.t = time.Time{}
		return nil
	}

	var ok bool
	t.t, ok = src.(time.Time)
	if !ok {
		return fmt.Errorf("can not convert %v to time.Time", src)
	}
	return nil
}

func (t Date) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *Date) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	t.t, err = time.Parse("2006-01-02", s)
	return err
}

func (t Date) MarshalJSON() ([]byte, error) {
	if t.t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.t.Format("2006-01-02"))
}

func (t Date) Add(y, m, d int) Date {
	return Date{t.t.AddDate(y, m, d)}
}

func (t Date) SubDay(t2 Date) int {
	return int(t.t.Sub(t2.t).Hours() / 24)
}

func (t Date) String() string {
	return t.t.Format("2006-01-02")
}

func (t Date) ParseThaiFormat() string {
	return t.t.Format("02/01/2006")
}

func (t Date) IsZero() bool {
	return t.t.IsZero()
}

func (t Date) Time() time.Time {
	return t.t
}

func (t Date) Year() int {
	return t.t.Year()
}

func (t Date) Month() time.Month {
	return t.t.Month()
}

func (t Date) Day() int {
	return t.t.Day()
}
