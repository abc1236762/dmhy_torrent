package main

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"strconv"
	"html"
	"fmt"
)

type Options struct {
	SortIDs, TeamIDs []IntStrPair
	Orders           []StrStrPair
}

func GetOptions(inclSortIDs, inclTeamIDs,
inclOrders bool) (opts *Options, err error) {
	opts = new(Options)
	var resp *http.Response
	var bodyBytes []byte
	var body string
	var origOpts []StrStrPair
	var key int
	
	if resp, err = http.Get(DmhyAdvSearchUrl); err != nil {
		return
	}
	if bodyBytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	body = string(bodyBytes)
	
	var getOrigOpts = func(name string) (origOpts []StrStrPair) {
		var optsStr = regexp.MustCompile(`<select name="` +
			name + `".*?>(.*?)</select>`).FindStringSubmatch(body)[1]
		for _, opt := range regexp.MustCompile(
			`<option value="([\w\-]+)".*?>(.*?)</option>`).
			FindAllStringSubmatch(optsStr, -1) {
			origOpts = append(origOpts, StrStrPair{
				opt[1], html.UnescapeString(opt[2])})
		}
		return
	}
	
	if inclSortIDs {
		origOpts = getOrigOpts("sort_id")
		for _, pair := range origOpts {
			if key, err = strconv.Atoi(pair.Key); err != nil {
				return
			}
			opts.SortIDs = append(opts.SortIDs, IntStrPair{key, pair.Val})
		}
	}
	
	if inclTeamIDs {
		origOpts = getOrigOpts("team_id")
		for _, pair := range origOpts {
			if key, err = strconv.Atoi(pair.Key); err != nil {
				return
			}
			opts.TeamIDs = append(opts.TeamIDs, IntStrPair{key, pair.Val})
		}
	}
	
	if inclOrders {
		opts.Orders = getOrigOpts("order")
	}
	
	return opts, resp.Body.Close()
}

func (o *Options) Print() {
	if o.SortIDs != nil {
		fmt.Println("List of sort_id and the meaning: ")
		fmt.Printf("%-8s\t%s\n", "*sort_id", "*meaning")
		for _, pair := range o.SortIDs {
			fmt.Printf("%-8v\t%s\n", pair.Key, pair.Val)
		}
		fmt.Println()
	}
	if o.TeamIDs != nil {
		fmt.Println("List of team_id and the meaning: ")
		fmt.Printf("%-8s\t%s\n", "*team_id", "*meaning")
		for _, pair := range o.TeamIDs {
			fmt.Printf("%-8v\t%s\n", pair.Key, pair.Val)
		}
		fmt.Println()
	}
	if o.Orders != nil {
		fmt.Println("List of order and the meaning: ")
		fmt.Printf("%-8s\t%s\n", "*order", "*meaning")
		for _, pair := range o.Orders {
			fmt.Printf("%-8v\t%s\n", pair.Key, pair.Val)
		}
		fmt.Println()
	}
}
