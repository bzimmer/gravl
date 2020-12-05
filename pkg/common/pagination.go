package common

import (
	"context"
	"errors"
)

// Pagination provides guidance on how to paginate through resources
type Pagination struct {
	// Total of resources to query
	Total int
	// Start at this page
	Start int
	// Count of the number of resources to query per page
	Count int
}

// Paginator paginates through results
type Paginator interface {
	// Page returns the default page size
	Page() int
	// Count of the number of resources queried
	Count() int
	// Do the querying
	Do(ctx context.Context, start, count int) (int, error)
}

func Paginate(ctx context.Context, paginator Paginator, spec Pagination) error {
	var (
		start = spec.Start
		count = spec.Count
		total = spec.Total
	)
	// log.Debug().
	// 	Str("prepared", "pre").
	// 	Int("start", start).
	// 	Int("count", count).
	// 	Int("total", total).
	// 	Msg("paginate")
	if total < 0 {
		return errors.New("total less than zero")
	}
	if start <= 0 {
		start = 1
	}
	if count <= 0 {
		count = paginator.Page()
	}
	if total > 0 && total <= count {
		count = total
	}
	// if requesting only one page of data then optimize
	if start <= 1 && total < paginator.Page() {
		count = total
	}
	// log.Debug().
	// 	Str("prepared", "post").
	// 	Int("start", start).
	// 	Int("count", count).
	// 	Int("total", total).
	// 	Msg("paginate")
	return do(ctx, paginator, total, start, count)
}

func do(ctx context.Context, paginator Paginator, total, start, count int) error {
	for {
		// all := paginator.Count()
		// log.Debug().
		// 	Int("all", all).
		// 	Int("start", start).
		// 	Int("count", count).
		// 	Int("total", total).
		// 	Msg("do")
		n, err := paginator.Do(ctx, start, count)
		if err != nil {
			return err
		}
		all := paginator.Count()
		// log.Debug().
		// 	Int("n", n).
		// 	Int("all", all).
		// 	Int("start", start).
		// 	Int("count", count).
		// 	Int("total", total).
		// 	Msg("done")
		// Strava documentation says receiving fewer than requested results is a
		// possible scenario so break only if 0 results were returned or we have
		// enough to fulfill the request
		if n == 0 || all >= total {
			break
		}
		start++
	}
	return nil
}
