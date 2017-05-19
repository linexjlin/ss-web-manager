package main

import "os"
import "io/ioutil"
import "encoding/json"

type MailGunConf struct {
	Sender, Domain, ApiKey, PublicApiKey string
}

type Conf struct {
	SSPortStartWith int
	UserIdStartWith int
	DefaultTraffic  int64
	DefaultCycle    int
	SiteName        string
	SiteLink        string
	MailGun         MailGunConf
}

func loadConf(confPath string, conf *Conf) {
	f, err := os.Open(confPath)
	defer f.Close()
	checkError(err)
	jsdat, err := ioutil.ReadAll(f)
	checkError(err)
	err = json.Unmarshal(jsdat, conf)
	checkError(err)
}
