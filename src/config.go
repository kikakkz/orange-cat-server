package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"kkt.com/glog"
)

type ContentSpecCfg struct {
	Host          string `json:"host"`
	Start         string `json:"start"`
	End           string `json:"end"`
	CharSet       string `json:"charset"`
	ChapterPrefix string `json:"chapter_prefix"`
	PType         string `json:"ptype"`
}

type MysqlCfg struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Db       string `json:"db"`
}

type config struct {
	Mysql       MysqlCfg         `json:"mysql"`
	ContentSpec []ContentSpecCfg `json:"content_spec"`
}

var cfg config

func init() {
	glog.ToStderr(true)
}

func ConfigInitialize() {
	cfgFile := flag.String("c", "./server-config.json", "Set `config file`")
	flag.Parse()

	body, err := ioutil.ReadFile(*cfgFile)
	if nil != err {
		glog.Error(err)
		return
	}

	err = json.Unmarshal(body, &cfg)
	if nil != err {
		glog.Error(err)
	}
}
