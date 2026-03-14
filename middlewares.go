package filter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OptionsContextKey is the gin.Context key under which [QueryOptionsMiddlewares]
// stores the parsed [Options].
var (
	OptionsContextKey = "filter_options"
)

const (
	defaultLimit        = 10
	defaultSortingField = "id"
)

// Params holds the query-parameter bindings for pagination and sorting.
// It is used by the Gin middlewares to parse incoming request parameters.
type (
	Params struct {
		SortBy     string `form:"sort_by" type:"string"`  // SortBy is the field name to sort by.
		Descending bool   `form:"descending" type:"bool"` // Descending enables descending sort order.
		Page       int    `form:"page" type:"int"`        // Page is the 1-based page number.
		Limit      int    `form:"limit" type:"int"`       // Limit is the maximum number of results per page.
	}
)

// QueryOptionsMiddlewares returns a [gin.HandlerFunc] that parses pagination and
// sorting parameters from the request query string. The resulting [Options] are
// stored in gin.Context under [OptionsContextKey].
//
// Supported query parameters: sort_by, descending, page, limit.
// Default values: limit=10, page=1, sort_by="id".
func QueryOptionsMiddlewares() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var filterParams Params
		if err := ctx.ShouldBind(&filterParams); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if filterParams.Limit == 0 {
			filterParams.Limit = defaultLimit
		}
		if filterParams.Page == 0 {
			filterParams.Page = 1
		}
		if filterParams.SortBy == "" {
			filterParams.SortBy = defaultSortingField
		}
		newOptions := NewOptions(filterParams.Limit, filterParams.Page, filterParams.SortBy, filterParams.Descending)
		ctx.Set(OptionsContextKey, newOptions)
		ctx.Next()
	}
}

// SingleQueryOptionsMiddlewares returns a [gin.HandlerFunc] that parses
// query parameters and passes the resulting [Options] directly to next.
// Default values: limit=10, page=1, sort_by="id".
func SingleQueryOptionsMiddlewares(next func(c *gin.Context, option Options)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var filterParams Params
		if err := ctx.ShouldBind(&filterParams); err != nil {
			return
		}
		if filterParams.Limit == 0 {
			filterParams.Limit = defaultLimit
		}
		if filterParams.Page == 0 {
			filterParams.Page = 1
		}
		if filterParams.SortBy == "" {
			filterParams.SortBy = defaultSortingField
		}
		newOptions := NewOptions(filterParams.Limit, filterParams.Page, filterParams.SortBy, filterParams.Descending)
		next(ctx, newOptions)
	}
}

// SingleQueryOptionsMiddlewaresWithDefaults is like
// [SingleQueryOptionsMiddlewares] but uses the provided defaults instead of
// the built-in defaults when a query parameter is missing.
func SingleQueryOptionsMiddlewaresWithDefaults(next func(c *gin.Context, option Options), defaults *Params) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var filterParams Params
		_, descendingProvided := ctx.GetQuery("descending")
		if err := ctx.ShouldBind(&filterParams); err != nil {
			return
		}

		if filterParams.Limit == 0 {
			filterParams.Limit = defaults.Limit
		}
		if filterParams.Page == 0 {
			filterParams.Page = defaults.Page
		}
		if filterParams.SortBy == "" {
			filterParams.SortBy = defaults.SortBy
		}
		if !descendingProvided {
			filterParams.Descending = defaults.Descending
		}
		newOptions := NewOptions(filterParams.Limit, filterParams.Page, filterParams.SortBy, filterParams.Descending)
		next(ctx, newOptions)
	}
}
