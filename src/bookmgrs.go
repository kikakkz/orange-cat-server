package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	_ "kkt.com/glog"
	"net/http"
	"strconv"
)

var mgr *BookMgr

type booksGetP struct {
	pageIndex int
	clazz     string
	gender    string
	finished  bool
	action    string
	key       string
}

type booksListResp struct {
	Books []*Book `json:"books"`
}

type booksSearchResp struct {
	TotalCount int     `json:"total_count"`
	Books      []*Book `json:"books"`
}

func createBookMgr() *BookMgr {
	mgr, _ := NewBookMgr()
	return mgr
}

func queryBooksList(clazz string, gender string, finished bool, curPage int) (*booksListResp, error) {
	books, err := mgr.QueryBooksList(clazz, gender, finished, curPage)
	if nil != err {
		return nil, err
	}
	var resp = booksListResp{Books: books}
	return &resp, nil
}

func queryBooksInfo(clazz string, gender string, finished bool) (*BooksInfo, error) {
	return mgr.QueryBooksInfo(clazz, gender, finished)
}

func searchBooks(clazz string, curPage int) (*booksSearchResp, error) {
	books, count, err := mgr.SearchBooks(clazz, curPage)
	if nil != err {
		return nil, err
	}
	var resp = booksSearchResp{Books: books, TotalCount: count}
	return &resp, nil
}

func queryBooks(p booksGetP) (interface{}, error) {
	if "l" == p.action {
		return queryBooksList(p.clazz, p.gender, p.finished, int(p.pageIndex))
	} else if "c" == p.action {
		return queryBooksInfo(p.clazz, p.gender, p.finished)
	} else if "s" == p.action {
		return searchBooks(p.clazz, int(p.pageIndex))
	}
	return nil, errors.New("Invalid action")
}

func booksGet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if nil != err {
		Response(w, -1, err.Error(), nil)
		return
	}

	var reqP booksGetP
	var pageIndex int64 = -1

	action := "c"
	if 0 < len(r.Form["a"]) {
		action = r.Form["a"][0]
	}
	reqP.action = action

	if 0 < len(r.Form["p"]) {
		pageIndex, err = strconv.ParseInt(r.Form["p"][0], 10, 32)
		if nil != err {
			Response(w, -2, err.Error(), nil)
			return
		}
		pageIndex -= 1
	}
	reqP.pageIndex = int(pageIndex)

	clazz := ""
	if 0 < len(r.Form["c"]) {
		clazz = r.Form["c"][0]
	}
	reqP.clazz = clazz

	finished := false
	if 0 < len(r.Form["f"]) {
		finished = r.Form["f"][0] == "true"
	}
	reqP.finished = finished

	gender := "default"
	if 0 < len(r.Form["g"]) {
		gender = r.Form["g"][0]
	}
	reqP.gender = gender

	resp, err := queryBooks(reqP)
	if nil != err {
		Response(w, -3, err.Error(), nil)
		return
	}

	Response(w, 0, "", resp)
}

func booksPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		Response(w, -1, err.Error(), nil)
		return
	}

	var reqP apiPostP
	err = json.Unmarshal(body, &reqP)
	if nil != err {
		Response(w, -2, err.Error(), nil)
		return
	}

	if "set" == reqP.Action {
		err = mgr.SetBooks(reqP.Key, reqP.Body)
	}

	if nil != err {
		Response(w, -3, err.Error(), nil)
		return
	}

	Response(w, 0, "", nil)
}

func BookMgrsProc(w http.ResponseWriter, r *http.Request) {
	if nil == mgr {
		mgr = createBookMgr()
	}

	switch r.Method {
	case "GET":
		booksGet(w, r)
	case "POST":
		booksPost(w, r)
	}
}
