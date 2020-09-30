package paginator

import (
	"errors"
	"math"
)

// DefaultMaxPerPage default number of records per page
const DefaultMaxPerPage = 10

// ErrNoPrevPage current page is first page
var ErrNoPrevPage = errors.New("no previous page")

// ErrNoNextPage current page is last page
var ErrNoNextPage = errors.New("no next page")

// Adapter any adapter must implement this interface
type Adapter interface {
	Nums() (int64, error)
	Slice(offset, length int, data interface{}) error
}

// Paginator structure
type Paginator struct {
	adapter    Adapter
	maxPerPage int
	page       int
	nums       int64
}

// New paginator constructor
func New(adapter Adapter, maxPerPage int) Paginator {
	if maxPerPage <= 0 {
		maxPerPage = DefaultMaxPerPage
	}

	return Paginator{
		adapter:    adapter,
		maxPerPage: maxPerPage,
		page:       1,
		nums:       -1,
	}
}

// SetPage set current page
func (p *Paginator) SetPage(page int) {
	if page <= 0 {
		page = 1
	}

	p.page = page
}

// Page returns current page
func (p Paginator) Page() (int, error) {
	pn, err := p.PageNums()
	if err != nil {
		return 0, err
	}

	if p.page > pn {
		return pn, nil
	}

	return p.page, nil
}

// Results stores the current page results into data argument which must be a pointer to a slice.
func (p Paginator) Results(data interface{}) error {
	var offset int
	page, err := p.Page()
	if err != nil {
		return err
	}

	if page > 1 {
		offset = (page - 1) * p.maxPerPage
	}

	return p.adapter.Slice(offset, p.maxPerPage, data)
}

// Nums returns the total number of records
func (p *Paginator) Nums() (int64, error) {
	var err error
	if p.nums == -1 {
		p.nums, err = p.adapter.Nums()
		if err != nil {
			return 0, err
		}
	}

	return p.nums, nil
}

// HasPages returns true if there is more than one page
func (p Paginator) HasPages() (bool, error) {
	n, err := p.Nums()
	if err != nil {
		return false, err
	}

	return n > int64(p.maxPerPage), nil
}

// HasNext returns true if current page is not the last page
func (p Paginator) HasNext() (bool, error) {
	pn, err := p.PageNums()
	if err != nil {
		return false, err
	}

	page, err := p.Page()
	if err != nil {
		return false, err
	}

	return page < pn, nil
}

// PrevPage returns previous page number or ErrNoPrevPage if current page is first page
func (p Paginator) PrevPage() (int, error) {
	hp, err := p.HasPrev()
	if err != nil {
		return 0, nil
	}

	if !hp {
		return 0, ErrNoPrevPage
	}

	page, err := p.Page()
	if err != nil {
		return 0, err
	}

	return page - 1, nil
}

// NextPage returns next page number or ErrNoNextPage if current page is last page
func (p Paginator) NextPage() (int, error) {
	hn, err := p.HasNext()
	if err != nil {
		return 0, err
	}

	if !hn {
		return 0, ErrNoNextPage
	}

	page, err := p.Page()
	if err != nil {
		return 0, err
	}

	return page, nil
}

// HasPrev returns true if current page is not the first page
func (p Paginator) HasPrev() (bool, error) {
	page, err := p.Page()
	if err != nil {
		return false, err
	}

	return page > 1, nil
}

// PageNums returns the total number of pages
func (p Paginator) PageNums() (int, error) {
	n, err := p.Nums()
	if err != nil {
		return 0, err
	}

	n = int64(math.Ceil(float64(n) / float64(p.maxPerPage)))
	if n == 0 {
		n = 1
	}

	return int(n), nil
}
