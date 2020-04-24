package main

type apiPostP struct {
	Action string      `json:"action"`
	Key    string      `json:"key"`
	Body   interface{} `json:"body"`
}
