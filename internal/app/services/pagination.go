package services

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage    int32 = 1
	DefaultPerPage int32 = 25
	MaxPerPage     int32 = 1000
)

type Pagination struct {
	Page    int32
	PerPage int32
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

func parseQueryParam(r *http.Request, key string, defaultValue int32) int32 {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue
	}

	value, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return defaultValue
	}

	return int32(value)
}

func (p *Pagination) Limit() int32 {
	return p.PerPage
}

func (p *Pagination) Offset() int32 {
	return (p.Page - 1) * p.PerPage
}
