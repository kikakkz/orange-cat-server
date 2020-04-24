package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	_ "kkt.com/glog"
	"net/http"
)

type bookGetP struct {
	action string
	id     string
	name   string
	author string
}

type bookChaptersP struct {
	Chapters []*Chapter `json:"chapters"`
}

func queryBookChapters(w http.ResponseWriter, p bookGetP, mgr *BookMgr) (*bookChaptersP, error) {
	chapters, err := mgr.GetBookChapters(p.name, p.author)
	if nil != err {
		return nil, err
	}

	var resp = bookChaptersP{Chapters: chapters}
	return &resp, nil
}

func queryBook(w http.ResponseWriter, p bookGetP) error {
	mgr, _ := NewBookMgr()
	var err error
	var resp interface{}

	switch p.action {
	case "ch":
		resp, err = queryBookChapters(w, p, mgr)
	}

	if nil != err {
		return err
	}

	Response(w, 0, "", resp)
	return nil
}

func bookGet(w http.ResponseWriter, r *http.Request) {
	var p bookGetP

	err := r.ParseForm()
	if nil != err {
		Response(w, -1, err.Error(), nil)
		return
	}

	var action = "ch"
	if 0 < len(r.Form["a"]) {
		action = r.Form["a"][0]
	}
	p.action = action

	var id = ""
	if 0 < len(r.Form["id"]) {
		id = r.Form["id"][0]
	}
	p.id = id

	var name = ""
	if 0 < len(r.Form["n"]) {
		name = r.Form["n"][0]
	}
	p.name = name

	var author = ""
	if 0 < len(r.Form["au"]) {
		author = r.Form["au"][0]
	}
	p.author = author

	if author == "" || name == "" || id == "" {
		Response(w, -2, "Invalid parameter", nil)
		return
	}

	err = queryBook(w, p)
	if nil != err {
		Response(w, -3, err.Error(), nil)
	}
}

type bookPostRespP struct {
	Book *Book `json:"book"`
}

func operateBook(p apiPostP) (*bookPostRespP, error) {
	mgr, _ := NewBookMgr()
	var book *Book
	var err error

	if "" == p.Action || "" == p.Key {
		return nil, errors.New("Invalid parameter")
	}

	switch p.Key {
	case "read":
		book, err = mgr.addBookRead(p.Body)
	}

	if nil != err {
		return nil, err
	}

	resp := bookPostRespP{Book: book}
	return &resp, nil
}

func bookPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		Response(w, -1, err.Error(), nil)
		return
	}

	var p apiPostP
	err = json.Unmarshal(body, &p)
	if nil != err {
		Response(w, -2, err.Error(), nil)
		return
	}

	resp, err := operateBook(p)
	if nil != err {
		Response(w, -3, err.Error(), nil)
		return
	}
	Response(w, 0, "", resp)
}

func BookProc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		bookGet(w, r)
	case "POST":
		bookPost(w, r)
	}
}
