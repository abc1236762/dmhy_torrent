package main

import (
	"flag"
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
)

const (
	DmhyUrl          = `https://share.dmhy.org`
	DmhyQueryUrl     = DmhyUrl + `/topics/list?keyword=%s&sort_id=%d&team_id=%d&user_id=%d&order=%s&page=%d`
	DmhyAdvSearchUrl = DmhyUrl + `/topics/advanced-search`
)

func main() {
	var fs *flag.FlagSet
	var err error
	
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Usage: <command> [arguments ... ]")
		fmt.Fprintln(os.Stderr, "  Use <command> -h or "+
			"<command> -help to get usage of the command.")
		fmt.Fprintln(os.Stderr, "Commands: ")
		fmt.Fprintln(os.Stderr, "  list\tList all sort_id, team_id, and order.")
		fmt.Fprintln(os.Stderr, "  get\tGet torrents and download them.")
		return
	} else {
		
		if os.Args[1] == "list" {
			fs = flag.NewFlagSet("command list", flag.ContinueOnError)
			var s = fs.Bool("s", false, "List all sort_id and its meaning.")
			var t = fs.Bool("t", false, "List all team_id and its meaning.")
			var o = fs.Bool("o", false, "List all order and its meaning.")
			if err = fs.Parse(os.Args[2:]); err != nil {
				if err != flag.ErrHelp {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				return
			}
			
			var opts *Options
			if opts, err = GetOptions(*s, *t, *o); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
			opts.Print()
			
		} else if os.Args[1] == "get" {
			
			fs = flag.NewFlagSet("command get", flag.ContinueOnError)
			var k = fs.String("k", "", "Set `keyword`.")
			var s = fs.Uint("s", 0, "Set `sort_id`.")
			var t = fs.Uint("t", 0, "Set `team_id`.")
			var u = fs.Uint("u", 0, "Set `user_id`.")
			var o = fs.String("o", "", "Set `order`.")
			var r = fs.String("r", "", "Set recent with "+
				"`<num>+<unit>` format, <num> is a positive number "+
				"and <unit> can be `day`, `week`, `month`, or `year`")
			var p = fs.String("p", "1-1", "Set page with `<begin>-<end>` "+
				"format, <begin> and <end> are positive numbers.")
			if err = fs.Parse(os.Args[2:]); err != nil {
				if err != flag.ErrHelp {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				return
			}
			
			var query *Query
			var errs []error
			if query, errs = MakeQuery(*k, *o, *r,
				*p, *s, *t, *u); len(errs) > 0 {
				for _, err := range errs {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				return
			}
			
			for _, item := range query.Items {
				if err = item.Download(); err != nil {
					errs = append(errs, err)
				}
			}
			for _, err := range errs {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			
		} else {
			fmt.Fprintf(os.Stderr, "Command %s is not supported.", os.Args[1])
		}
	}
}

func HttpGetBodyStr(u string) (body string, err error) {
	var resp *http.Response
	var bodyBytes []byte
	
	if resp, err = http.Get(u); err != nil {
		return
	}
	if bodyBytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	return string(bodyBytes), resp.Body.Close()
}

func FixFilename(name string) string {
	return regexp.MustCompile(`[\\/:*?"<>|]`).ReplaceAllString(name, "_")
}