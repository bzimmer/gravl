package wta

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
)

var photosRE = regexp.MustCompile(`([0-9]+)`)

// ReportsService .
type ReportsService service

func query(author string) *url.URL {
	v := url.Values{}
	v.Add("author", author)
	v.Add("b_size", "100")
	// v.Add("b_start:int", "0")
	v.Add("filter", "Search")
	v.Add("hiketypes:list", "day-hike")
	v.Add("hiketypes:list", "multi-night-backpack")
	v.Add("hiketypes:list", "overnight")
	v.Add("hiketypes:list", "snowshoe-xc-ski")
	v.Add("month", "all")
	v.Add("subregion", "all")

	// parsing a constant, if this fails we have other issues
	u, _ := url.Parse(baseURL)
	u.RawQuery = v.Encode()
	return u
}

func newCollector(client *http.Client) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("wta.org", "www.wta.org"),
	)
	c.SetClient(client)
	return c
}

// TripReports .
func (s *ReportsService) TripReports(ctx context.Context, reporter string) ([]*TripReport, error) {
	var visitError error
	reports := make([]*TripReport, 0)

	q := query(reporter).String()
	c := newCollector(s.client.client)
	c.OnError(func(r *colly.Response, err error) {
		log.Warn().
			Err(err).
			Str("url", r.Request.URL.String()).
			Msg("tripreports")
		visitError = err
	})

	c.OnHTML("div[class=item-row]", func(e *colly.HTMLElement) {
		tr := &TripReport{
			Title:  e.ChildText(".listitem-title"),
			Region: e.ChildText("span[class=region]"),
		}

		creator := strings.Split(e.ChildTexts("div[class=CreatorInfo]")[0], "\n")

		tr.Report = e.ChildAttr(".listitem-title", "href")
		tr.Reporter = creator[0]
		txt := e.ChildText(".UpvoteCount")
		if txt != "" {
			vote, err := strconv.Atoi(txt)
			if err == nil {
				tr.Votes = vote
			}
		}
		txt = e.ChildText(".media-indicator")
		if txt != "" {
			n := photosRE.FindString(txt)
			photos, err := strconv.Atoi(n)
			if err == nil {
				tr.Photos = photos
			}
		}
		attr := e.ChildAttr(".elapsed-time", "title")
		if attr != "" {
			t, _ := time.Parse("Jan 02, 2006", attr)
			tr.HikeDate = t
		}

		reports = append(reports, tr)
	})

	defer func(start time.Time) {
		log.Debug().
			Str("url", q).
			Str("op", "reports").
			Dur("elapsed", time.Since(start)).
			Msg("GetTripReports")
	}(time.Now())

	err := c.Visit(q)
	if err != nil {
		return nil, err
	}
	if visitError != nil {
		return nil, visitError
	}
	return reports, nil
}
