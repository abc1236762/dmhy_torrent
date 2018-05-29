package main

import (
	"fmt"
	"errors"
	"regexp"
)

type Query struct {
	Keyword, Order, Recent string
	SortID, TeamID         int
	Page                   *Page
}

func (q *Query) Check() (errs []error) {
	var opts *Options
	var err error
	
	if opts, err = GetOptions(true, true, true); err != nil {
		return []error{err}
	}
	
	if !IntStrPairsHasKey(opts.SortIDs, q.SortID) {
		errs = append(errs,
			fmt.Errorf("query: sort_id %d is not found", q.SortID))
	}
	if !IntStrPairsHasKey(opts.TeamIDs, q.TeamID) {
		errs = append(errs,
			fmt.Errorf("query: team_id %d is not found", q.TeamID))
	}
	if len(q.Order) > 0 {
		if !StrStrPairsHasKey(opts.Orders, q.Order) {
			errs = append(errs,
				fmt.Errorf("query: order %s is not found", q.Order))
		}
	}
	if q.Page.Begin <= 0 || q.Page.End <= 0 {
		errs = append(errs, errors.New("query: <begin> or " +
			"<end> of page should be a positive number"))
	}
	if q.Page.Begin > q.Page.End {
		errs = append(errs, errors.New(
			"query: <begin> of page should not be more than <end> of page"))
	}
	if len(q.Recent) > 0 {
		if !regexp.MustCompile(
			`^\d+\+(day|week|month|year)$`).MatchString(q.Recent) {
			errs = append(errs, fmt.Errorf(
				"query: recent %s has wrong format", q.Recent))
		} else if q.Recent[0] - '0' == 0 {
			errs = append(errs, errors.New(
				"query: <num> of recent should be a positive number"))
		}
	}
	
	return
}

func (q *Query) GetUrls() (urls []string) {
	for i := q.Page.Begin; i <= q.Page.End; i++ {
		var url = fmt.Sprintf(DmhyQueryUrl,
			q.Keyword, q.SortID, q.TeamID, q.Order, i)
		if len(q.Recent) > 0 {
			url += "&recent=" + q.Recent
		}
		urls = append(urls, url)
	}
	return
}
