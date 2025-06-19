package api

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

type Query struct {
	s   string
	raw string
}

func NewQuery(s string) Query {
	s = strings.TrimSpace(s)
	return Query{queryReplacer.Replace(s), s}
}

var queryReplacer = strings.NewReplacer("%", "\\%", "_", "\\_")

func (q *Query) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*q = NewQuery(s)
	return nil
}

func (q Query) String() string {
	return q.s
}

func (q Query) Value() (driver.Value, error) {
	return "%" + q.s + "%", nil
}

func (q Query) IsZero() bool {
	return q.s == ""
}

func (q Query) Raw() string {
	return q.raw
}
