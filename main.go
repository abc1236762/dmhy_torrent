package main

import (
	"flag"
	"os"
	"fmt"
)

const (
	DmhyUrl          = `https://share.dmhy.org`
	DmhyQueryUrl     = DmhyUrl + `/topics/list?keyword=%s&sort_id=%d&team_id=%d&order=%s&page=%d`
	DmhyAdvSearchUrl = DmhyUrl + `/topics/advanced-search`
)

func main() {
	var fs = new(flag.FlagSet)
	var err error
	if len(os.Args) <= 1 {
		// TODO: Usage of command.
		return
	} else {
		if os.Args[1] == "list" {
			var s = fs.Bool("s", false, "List all sort_id and its meaning.")
			var t = fs.Bool("t", false, "List all team_id and its meaning.")
			var o = fs.Bool("o", false, "List all order and its meaning.")
			if err = fs.Parse(os.Args[2:]); err != nil {
				if err != flag.ErrHelp {
					fmt.Println(err.Error())
				}
				return
			}
			
			var opts *Options
			if opts, err = GetOptions(*s, *t, *o); err != nil {
				fmt.Println(err.Error())
				return
			}
			opts.Print()
		} else if os.Args[1] == "get" {
			var k = fs.String("k", "", "Set `keyword`.")
			var s = fs.Int("s", 0, "Set `sort_id`.")
			var t = fs.Int("t", 0, "Set `team_id`.")
			var o = fs.String("o", "", "Set `order`.")
			var r = fs.String("r", "", "Set recent with " +
				"`<num>+<unit>` format, <num> is a positive number " +
				"and <unit> can be `day`, `week`, `month`, or `year`")
			var p = fs.String("p", "1-1", "Set page with `<begin>-<end>` " +
				"format, <begin> and <end> are positive numbers.")
			if err = fs.Parse(os.Args[2:]); err != nil {
				if err != flag.ErrHelp {
					fmt.Println(err.Error())
				}
				return
			}
			
			var page *Page
			if page, err = MakePage(*p); err != nil {
				fmt.Println(err.Error())
				return
			}
			var query = &Query{*k, *o, *r, *s, *t, page}
			var errs = query.Check()
			if errs != nil {
				for _, err := range errs {
					fmt.Println(err.Error())
				}
				return
			}
			for _, url := range query.GetUrls() {
				fmt.Println(url)
			}
		}
	}
}
