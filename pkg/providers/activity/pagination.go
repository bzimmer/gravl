package activity

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

// Pagination specifies how to paginate through resources
type Pagination struct {
	// Total number of resources to query
	Total int
	// Start querying at this page
	Start int
	// Count of the number of resources to query per page
	Count int
}

// Paginator paginates through results
type Paginator interface {
	// PageSize returns the number of resources to query per request
	PageSize() int
	// Count of the aggregate total of resources queried
	Count() int
	// Do executes the query using the pagination specification returning
	// the number of resources returned in this request or an error
	Do(ctx context.Context, spec Pagination) (int, error)
}

func Paginate(ctx context.Context, paginator Paginator, spec Pagination) error {
	var (
		start = spec.Start
		count = spec.Count
		total = spec.Total
	)
	if total < 0 {
		return errors.New("total less than zero")
	}
	if start <= 0 {
		start = 1
	}
	if count <= 0 {
		count = paginator.PageSize()
	}
	if total > 0 {
		if total <= count {
			count = total
		}
		// if requesting only one page of data then optimize
		if start <= 1 && total < paginator.PageSize() {
			count = total
		}
	}
	return do(ctx, paginator, Pagination{Total: total, Start: start, Count: count})
}

func do(ctx context.Context, paginator Paginator, spec Pagination) error {
	log.Info().
		Int("n", 0).
		Int("all", 0).
		Int("start", 0).
		Int("count", spec.Count).
		Int("total", spec.Total).
		Msg("do")
	for {
		n, err := paginator.Do(ctx, spec)
		if err != nil {
			return err
		}
		all := paginator.Count()
		// @warning(bzimmer)
		// the `spec.Count` value must be consistent throughout the entire pagination
		log.Info().
			Int("n", n).
			Int("all", all).
			Int("start", spec.Start).
			Int("count", spec.Count).
			Int("total", spec.Total).
			Msg("do")
		// fewer than requested results is a possible scenario so break only if
		//  0 results were returned or we have enough to fulfill the request
		if n == 0 {
			break
		}
		if spec.Total > 0 && all >= spec.Total {
			break
		}
		spec.Start++
	}
	return nil
}
