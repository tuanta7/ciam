package middleware

import (
	"context"
	"net/http"
	"strconv"
)

const (
	defaultPage          = 1
	defaultPageSize      = 10
	defaultTotalRequired = false

	ctxKeyPaginationParams = "pagination_params"
)

type paginationParams struct {
	Page          int32
	PageSize      int32
	TotalRequired bool
}

// Pagination middleware is used to handle pagination parameters when using offset-based pagination.
// It extracts pagination parameters from the request and sets them in the context for use in handlers.
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set default values for pagination parameters
		page := defaultPage
		if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
			page = p
		}

		pageSize := defaultPageSize
		if ps, err := strconv.Atoi(r.URL.Query().Get("page_size")); err == nil && ps > 0 {
			pageSize = ps
		}

		totalRequired := r.URL.Query().Get("total_required") == "true"

		ctx := context.WithValue(r.Context(), ctxKeyPaginationParams, paginationParams{
			Page:          int32(page),
			PageSize:      int32(pageSize),
			TotalRequired: totalRequired,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPaginationParams is a helper function to retrieve pagination params from the context.
func GetPaginationParams(ctx context.Context) (page int32, pageSize int32, totalRequired bool) {
	if params, ok := ctx.Value(ctxKeyPaginationParams).(paginationParams); ok {
		return params.Page, params.PageSize, params.TotalRequired
	}

	return defaultPage, defaultPageSize, defaultTotalRequired
}
