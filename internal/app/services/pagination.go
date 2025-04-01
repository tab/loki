package services

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage    uint64 = 1
	DefaultPerPage uint64 = 25
	MaxPerPage     uint64 = 1000
)

type Pagination struct {
	Page    uint64
	PerPage uint64
}

func NewPagination(r *http.Request) *Pagination {
	page := parseQueryParam(r, "page", DefaultPage)
	per := parseQueryParam(r, "per", DefaultPerPage)

	if page < 1 {
		page = DefaultPage
	}

	if per < 1 {
		per = DefaultPerPage
	}

	if per > MaxPerPage {
		per = MaxPerPage
	}

	return &Pagination{
		Page:    page,
		PerPage: per,
	}
}

func parseQueryParam(r *http.Request, key string, defaultValue uint64) uint64 {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue
	}

	value, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

func (p *Pagination) Limit() uint64 {
	return p.PerPage
}

func (p *Pagination) Offset() uint64 {
	return (p.Page - 1) * p.PerPage
}
