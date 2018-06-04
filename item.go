package main

import (
	"time"
	"regexp"
	"net/http"
	"io/ioutil"
	"html"
	"fmt"
	"os"
)

type Item struct {
	ID, Title, Team, TorrentUrl string
	Time                        time.Time
}

func NewItem(u string) (item *Item, err error) {
	var body string
	if body, err = HttpGetBodyStr(u); err != nil {
		return
	}
	
	item = new(Item)
	item.ID = regexp.MustCompile(`view/(\d+)_`).FindStringSubmatch(u)[1]
	var match []string
	
	match = regexp.MustCompile(
		`<a href="(//dl\.dmhy\.org/\d+/\d+/\d+/\w+\.torrent)">(.*?)</a>`).
		FindStringSubmatch(body)
	item.Title = html.UnescapeString(match[2])
	item.TorrentUrl = "https:" + match[1]
	
	match = regexp.MustCompile(
		`<span>(\d+/\d+/\d+ \d+:\d+)</span>`).FindStringSubmatch(body)
	if item.Time, err = time.Parse("2006/01/02 15:04 -0700",
		match[1]+" +0800"); err != nil {
		return
	}
	
	if match = regexp.MustCompile(
		`<a href="/topics/list/team_id/\d+">(.*?)</a>`).
		FindStringSubmatch(body); len(match) > 0 {
		item.Team = html.UnescapeString(match[1])
	}
	
	return
}

func (i *Item) Download() (err error) {
	// var team = " " + i.Team
	// if len(i.Team) == 0 {
	// 	team = ""
	// }
	// var path = "./" + FixFilename(fmt.Sprintf("%s%s - %s.torrent",
	// 	i.Time.Format("0601021504"), team, i.Title)) + ""
	var team = i.Team + " - "
	if len(i.Team) == 0 {
		team = ""
	}
	var path = "./" + FixFilename(fmt.Sprintf("%s%s.%s.torrent",
		team, i.Title, i.ID))
	fmt.Printf("Saving %s\n\n", path)
	
	var resp *http.Response
	var bodyBytes []byte
	var fileTime time.Time
	if resp, err = http.Get(i.TorrentUrl); err != nil {
		return
	}
	if bodyBytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if err = ioutil.WriteFile(path, bodyBytes, 0666); err != nil {
		return
	}
	if fileTime, err = time.Parse(time.RFC1123,
		resp.Header.Get("Last-Modified")); err != nil {
		return
	}
	if err = os.Chtimes(path, fileTime, fileTime); err != nil {
		return
	}
	
	return
}
