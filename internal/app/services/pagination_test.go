package services

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPagination(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected *Pagination
	}{
		{
			name: "Empty params",
			path: "/",
			expected: &Pagination{
				Page:    1,
				PerPage: 25,
			},
		},
		{
			name: "Valid params",
			path: "/?page=10&per=20",
			expected: &Pagination{
				Page:    10,
				PerPage: 20,
			},
		},
		{
			name: "Negative page",
			path: "/?page=-1",
			expected: &Pagination{
				Page:    1,
				PerPage: 25,
			},
		},
		{
			name: "Zero per",
			path: "/?per=0",
			expected: &Pagination{
				Page:    1,
				PerPage: 25,
			},
		},
		{
			name: "PerPage exceeds MaxPerPagePage",
			path: "/?per=5000",
			expected: &Pagination{
				Page:    1,
				PerPage: 1000,
			},
		},
		{
			name: "Invalid number format",
			path: "/?page=abc&per=xyz",
			expected: &Pagination{
				Page:    1,
				PerPage: 25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", tt.path, nil)
			paginator := NewPagination(request)
			assert.Equal(t, tt.expected.Page, paginator.Page)
			assert.Equal(t, tt.expected.PerPage, paginator.PerPage)
		})
	}
}

func Test_Pagination_Offset(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected uint64
	}{
		{
			name:     "Empty params",
			path:     "/",
			expected: 0,
		},
		{
			name:     "Valid params",
			path:     "/?page=10&per=20",
			expected: 180,
		},
		{
			name:     "Negative page",
			path:     "/?page=-1",
			expected: 0,
		},
		{
			name:     "Zero per",
			path:     "/?per=0",
			expected: 0,
		},
		{
			name:     "PerPage exceeds MaxPerPagePage",
			path:     "/?per=5000",
			expected: 0,
		},
		{
			name:     "Invalid number format",
			path:     "/?page=abc&per=xyz",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", tt.path, nil)
			paginator := NewPagination(request)
			assert.Equal(t, tt.expected, paginator.Offset())
		})
	}
}
