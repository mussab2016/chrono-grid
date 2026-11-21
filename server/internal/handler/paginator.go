package handler

import (
	"net/url"
	"strconv"
	"strings"
	"unisched/server/internal/dto"
)

func paginatorParams(query url.Values) *dto.PaginationQueryParams {
	// 1. Parsing Page Range (e.g., "0-15")
	pageParam := strings.TrimSpace(query.Get("page"))
	offset, limit := 0, 0 // 0 means fetch all by default if unspecified

	if pageParam != "" {
		parts := strings.Split(pageParam, "-")
		if len(parts) == 2 {
			pageIndex, err1 := strconv.Atoi(parts[0])
			pageSize, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil && pageSize > 0 {
				limit = pageSize
				offset = pageIndex * pageSize
			}
		}
	}

	// 2. Extracting & Cleaning optional params
	filter := strings.Trim(strings.TrimSpace(query.Get("filter")), `"`)
	sort := strings.Trim(strings.TrimSpace(query.Get("sort")), `"`)

	return &dto.PaginationQueryParams{
		Offset: offset,
		Limit:  limit,
		Filter: filter,
		SortBy: sort,
	}
}
