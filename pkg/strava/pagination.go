package strava

import (
	"errors"

	"github.com/rs/zerolog/log"
)

const (
	// PageSize of a default pagination request
	PageSize = 100
)

// Paginator paginates through results
type Paginator interface {
	Count() int
	Do(start, count int) (int, error)
}

func paginate(paginator Paginator, specs ...int) error {
	var start, count, total int
	switch len(specs) {
	case 0:
		total, start, count = 0, 1, PageSize
	case 1:
		total, start, count = specs[0], 1, PageSize
	case 2:
		total, start, count = specs[0], specs[1], PageSize
	case 3:
		total, start, count = specs[0], specs[1], specs[2]
	default:
		return errors.New("too many varargs")
	}
	log.Debug().
		Int("start", start).
		Int("count", count).
		Int("total", total).
		Ints("specs", specs).
		Msg("paginate")
	if total < 0 {
		return errors.New("total less than zero")
	}
	if total > 0 && total <= count {
		count = total
	}
	return do(paginator, total, start, count)
}

func do(paginator Paginator, total, start, count int) error {
	for {
		log.Debug().
			Int("start", start).
			Int("count", count).
			Int("total", total).
			Msg("do")
		n, err := paginator.Do(start, count)
		if err != nil {
			return err
		}
		all := paginator.Count()
		if n != count || all >= total {
			break
		}
		start++

		// The original implementation of pagination reset the count from `pageSize`
		// to the number of records required to fulfill the request if the remainder
		// was less than `pageSize`. This results in Strava returning the right number
		// of remaining records but they are duplicates from the first page! I was able
		// to reproduce this consistently. The Strava pagination document basically
		// reads as a best effort approach (eg ignore the result count and just keep
		// paging until no records are returned).
		if count > PageSize {
			// do not optimize count, Strava doesn't like
			count = PageSize
		} else if start <= 1 && total < PageSize {
			// unless it's the first pass through and you will not need further pagination
			count = total
		}
	}
	return nil
}
