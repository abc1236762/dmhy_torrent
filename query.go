package main

import (
	"fmt"
	"errors"
	"regexp"
	"strings"
	"strconv"
)

type Page struct {
	Begin, End uint64
}

type Query struct {
	Keyword, Order, Recent string
	SortID, TeamID, UserID uint
	Page                   *Page
	Items                  []*Item
}

func MakeQuery(k, o, r, p string, s, t, u uint) (query *Query, errs []error) {
	query = &Query{k, o, r, s, t, u, nil, nil}
	errs = append(errs, query.parsePage(p) ...)
	errs = append(errs, query.check() ...)
	if len(errs) > 0 {
		return
	}
	errs = append(errs, query.getItems() ...)
	return
}

func (q *Query) parsePage(p string) (errs []error) {
	var pageStr = strings.Split(p, "-")
	var err error
	var page = new(Page)
	
	if len(pageStr) != 2 {
		return []error{fmt.Errorf("query: page %s has wrong format", p)}
	}
	if page.Begin, err = strconv.ParseUint(pageStr[0], 10, 0); err != nil {
		errs = append(errs, errors.New("query: <begin> of page "+
			pageStr[0]+ " is not a positive number"))
	}
	if page.End, err = strconv.ParseUint(pageStr[1], 10, 0); err != nil {
		errs = append(errs, errors.New("query: <end> of page "+
			pageStr[1]+ " is not a positive number"))
	}
	
	if len(errs) == 0 {
		q.Page = page
	}
	return
}

func (q *Query) check() (errs []error) {
	var opts *Options
	var err error
	
	if opts, err = GetOptions(true, true, true); err != nil {
		return []error{err}
	}
	
	if !UintStrPairsHasKey(opts.SortIDs, q.SortID) {
		errs = append(errs,
			fmt.Errorf("query: sort_id %d is not found", q.SortID))
	}
	if !UintStrPairsHasKey(opts.TeamIDs, q.TeamID) {
		errs = append(errs,
			fmt.Errorf("query: team_id %d is not found", q.TeamID))
	}
	if len(q.Order) > 0 {
		if !StrStrPairsHasKey(opts.Orders, q.Order) {
			errs = append(errs,
				fmt.Errorf("query: order %s is not found", q.Order))
		}
	}
	if q.Page != nil {
		if q.Page.Begin <= 0 || q.Page.End <= 0 {
			errs = append(errs, errors.New("query: <begin> or "+
				"<end> of page should be a positive number"))
		}
		if q.Page.Begin > q.Page.End {
			errs = append(errs, errors.New(
				"query: <begin> of page should not be more than <end> of page"))
		}
	}
	
	if len(q.Recent) > 0 {
		if !regexp.MustCompile(
			`^\d+\+(day|week|month|year)$`).MatchString(q.Recent) {
			errs = append(errs, fmt.Errorf(
				"query: recent %s has wrong format", q.Recent))
		} else if q.Recent[0]-'0' == 0 {
			errs = append(errs, errors.New(
				"query: <num> of recent should be a positive number"))
		}
	}
	
	return
}

func (q *Query) getUrls() (urls []string) {
	for i := q.Page.Begin; i <= q.Page.End; i++ {
		var url = fmt.Sprintf(DmhyQueryUrl,
			q.Keyword, q.SortID, q.TeamID, q.UserID, q.Order, i)
		if len(q.Recent) > 0 {
			url += "&recent=" + q.Recent
		}
		urls = append(urls, url)
	}
	return
}

func (q *Query) getItems() (errs []error) {
	var urls = q.getUrls()
	var body string
	var item *Item
	var itemsMatch []string
	var err error
	for _, url := range urls {
		if body, err = HttpGetBodyStr(url); err != nil {
			errs = append(errs, err)
			continue
		}
		if itemsMatch = regexp.MustCompile(`/topics/view/.*?\.html`).
			FindAllString(body, -1); len(itemsMatch) == 0 {
			break
		}
		for _, sub := range itemsMatch {
			if item, err = NewItem(DmhyUrl + sub); err != nil {
				errs = append(errs, err)
				continue
			}
			fmt.Printf("Getted %s %v %s\n%s\n\n", item.Title,
				item.Time, item.Team, item.TorrentUrl)
			q.Items = append(q.Items, item)
		}
	}
	return
}
