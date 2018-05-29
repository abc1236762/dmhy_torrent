package main

import (
	"strings"
	"errors"
	"strconv"
)

type Page struct {
	Begin, End uint64
}

func MakePage(pageArg string) (page *Page, err error) {
	page = new(Page)
	var pageStr = strings.Split(pageArg, "-")
	if len(pageStr) != 2 {
		return nil, errors.New("page: format is wrong")
	}
	if page.Begin, err = strconv.ParseUint(pageStr[0], 10, 0); err != nil {
		return nil, errors.New("page: "+pageStr[0]+" is not a number")
	}
	if page.End, err = strconv.ParseUint(pageStr[1], 10, 0); err != nil {
		return nil, errors.New("page: "+pageStr[1]+" is not a number")
	}
	return
}