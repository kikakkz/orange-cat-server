package main

import (
	"bytes"
	"encoding/json"
	"kkt.com/glog"
	"net/http"
)

type resp struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Body  interface{} `json:"body,omitempty"`
}

func jsonMarshal(r interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(r)
	return buffer.Bytes(), err
}

func Response(w http.ResponseWriter, code int, err string, body interface{}) {
	var r = resp{Code: code, Error: err}
	if nil != body {
		r.Body = body
	}
	rJSON, e := jsonMarshal(r)
	if nil != e {
		glog.Error(e)
		return
	}
	w.Write(rJSON)
}
