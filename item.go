package main

import (
	"time"
	"regexp"
	"net/http"
	"io/ioutil"
	"html"
	"fmt"
)

type Item struct {
	Title, Team, TorrentUrl string
	Time                    *time.Time
}

func NewItem(url string) (item *Item, err error) {
	var body string
	if body, err = HttpGetBodyStr(url); err != nil {
		return
	}
	
	item = new(Item)
	item.Time = new(time.Time)
	var submatch []string
	
	submatch = regexp.MustCompile(
		`<a href="(//dl\.dmhy\.org/\d+/\d+/\d+/\w+\.torrent)">(.*?)</a>`).
		FindStringSubmatch(body)
	item.Title = html.UnescapeString(submatch[2])
	item.TorrentUrl = "https:" + submatch[1]
	
	submatch = regexp.MustCompile(
		`<span>(\d+/\d+/\d+ \d+:\d+)</span>`).FindStringSubmatch(body)
	if *item.Time, err = time.Parse("2006/01/02 15:04 -0700",
		submatch[1]+" +0800"); err != nil {
		return
	}
	
	if submatch = regexp.MustCompile(
		`<a href="/topics/list/team_id/\d+">(.*?)</a>`).
		FindStringSubmatch(body); len(submatch) > 0 {
			item.Team = html.UnescapeString(submatch[1])
	}
	
	return
}

func (i *Item) Download() (err error) {
	var team = " " + i.Team
	if len(i.Team) == 0 {
		team = ""
	}
	var path = "./" + FixFilename(fmt.Sprintf("%s%s - %s",
		i.Time.Format("0601021504"), team, i.Title)) + ".torrent"
	fmt.Printf("Saving %s\n\n", path)
	
	var resp *http.Response
	var bodyBytes []byte
	if resp, err = http.Get(i.TorrentUrl); err != nil {
		return
	}
	if bodyBytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if err = ioutil.WriteFile(path, bodyBytes, 0666); err != nil {
		return
	}
	
	return
}
