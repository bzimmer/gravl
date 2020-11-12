package strava

import (
	"errors"

	"github.com/rs/zerolog/log"
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
		total, start, count = 0, 1, pageSize
	case 1:
		total, start, count = specs[0], 1, pageSize
	case 2:
		total, start, count = specs[0], specs[1], pageSize
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
	return doPaginate(paginator, total, start, count)
}

func doPaginate(paginator Paginator, total, start, count int) error {
	log.Debug().
		Int("start", start).
		Int("count", count).
		Int("total", total).
		Msg("doPaginate")
	for {
		n, err := paginator.Do(start, count)
		if err != nil {
			return err
		}
		all := paginator.Count()
		if n != count || all >= total {
			break
		}
		start = start + 1
		if (total - all) < pageSize {
			count = total - all
		} else {
			count = pageSize
		}
	}
	return nil
}
