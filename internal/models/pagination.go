package models

const (
	DefaultPage     = 1
	DefaultLimit    = 10
	DefaultMaxLimit = 100
)

type Pagination struct {
	Limit     int
	Page      int
	SortBy    string
	SortOrder string
	MaxLimit  uint
}

func NewPagination(page, limit int) *Pagination {
	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}

func (p *Pagination) GetPage() int {
	if p.Page < 1 {
		return DefaultPage
	}
	return p.Page
}

func (p *Pagination) GetLimit() int {
	limit := DefaultLimit
	if p.Limit >= 1 {
		limit = p.Limit
	}
	if p.MaxLimit > 0 && limit > int(p.MaxLimit) {
		return int(p.MaxLimit)
	}
	return limit
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetSortBy() string {
	return p.SortBy
}

func (p *Pagination) GetSortOrder() string {
	if p.SortOrder == "" {
		return "asc"
	}
	return p.SortOrder
}

func (p *Pagination) WithMaxLimit(limit uint) *Pagination {
	if limit == 0 {
		p.MaxLimit = DefaultMaxLimit
	} else {
		p.MaxLimit = limit
	}
	return p
}

func (p *Pagination) WithSortBy(sortBy string) *Pagination {
	p.SortBy = sortBy
	return p
}

func (p *Pagination) WithSortOrder(sortOrder string) *Pagination {
	p.SortOrder = sortOrder
	return p
}
