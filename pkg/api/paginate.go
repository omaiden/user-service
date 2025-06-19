package api

import (
	"encoding/json"

	"github.com/acoshift/paginate"
)

const (
	DefaultPerPage = 25
	MaxPerPage     = 100
)

type Paginate struct {
	p *paginate.MovablePaginate
}

func (p *Paginate) init() {
	if p.p == nil {
		p.p = paginate.NewMovable(1, DefaultPerPage, 0)
	}
}

func (p *Paginate) MarshalJSON() ([]byte, error) {
	p.init()
	return json.Marshal(struct {
		Page    int64 `json:"page"`
		PerPage int64 `json:"perPage"`
		Next    bool  `json:"next"`
	}{
		p.p.Page(),
		p.p.PerPage(),
		p.p.CanNext(),
	})
}

func (p *Paginate) Count(f func() (int64, error)) error {
	p.init()
	cnt, err := f()
	if err != nil {
		return err
	}
	p.p.SetCount(cnt)
	return nil
}

func (p *Paginate) UnmarshalJSON(b []byte) error {
	var s struct {
		Page    int64 `json:"page"`
		PerPage int64 `json:"perPage"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	if s.PerPage <= 0 {
		s.PerPage = DefaultPerPage
	}
	if s.PerPage > MaxPerPage {
		s.PerPage = MaxPerPage
	}
	*p = Paginate{
		p: paginate.NewMovable(s.Page, s.PerPage, 0),
	}
	return nil
}

func (p *Paginate) CountOffset() int64 {
	p.init()
	return p.p.CountOffset()
}

func (p *Paginate) CountLimit() int64 {
	p.init()
	return p.p.CountLimit()
}

func (p *Paginate) Offset() int64 {
	p.init()
	return p.p.Offset()
}

func (p *Paginate) Limit() int64 {
	p.init()
	return p.p.Limit()
}

func (p *Paginate) Page() int64 {
	p.init()
	return p.p.Page()
}

func (p *Paginate) PerPage() int64 {
	p.init()
	return p.p.PerPage()
}
