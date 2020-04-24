package main

import (
	"kkt.com/glog"
	"net/http"
)

func serveBooks(w http.ResponseWriter, r *http.Request) {
	BookMgrsProc(w, r)
}

func serveBook(w http.ResponseWriter, r *http.Request) {
	BookProc(w, r)
}

func main() {
	ConfigInitialize()
	err := DBOpen(&cfg.Mysql)
	if nil != err {
		glog.Error(err)
		return
	}
	http.HandleFunc("/books", serveBooks)
	http.HandleFunc("/book", serveBook)
	for {
		http.ListenAndServe(":8999", nil)
	}
	DBClose()
}
