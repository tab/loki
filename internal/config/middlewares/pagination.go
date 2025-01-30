package middlewares

import (
	"context"
	"net/http"

	"loki/internal/app/services"
	"loki/pkg/logger"
)

type PaginationKey struct{}

type PaginationMiddleware interface {
	Paginate(next http.Handler) http.Handler
}

type paginationMiddleware struct {
	pagination services.Pagination
	log        *logger.Logger
}

func NewPaginationMiddleware(pagination services.Pagination, log *logger.Logger) PaginationMiddleware {
	return &paginationMiddleware{
		pagination: pagination,
		log:        log,
	}
}

func (p *paginationMiddleware) Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := withPagination(r.Context(), p.pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withPagination(ctx context.Context, pagination services.Pagination) context.Context {
	return context.WithValue(ctx, PaginationKey{}, pagination)
}

func CurrentPaginationFromContext(ctx context.Context) *services.Pagination {
	pagination, ok := ctx.Value(PaginationKey{}).(*services.Pagination)
	if !ok {
		return services.NewPagination(nil)
	}

	return pagination
}
