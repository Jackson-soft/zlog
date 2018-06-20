package main

import "zlog"

func main() {
	zlog.Infoln("df")
	zlog.WithFields(zlog.Fields{"fff": "dff", "vvv": 45.6}).Infoln("dfsdfa")
	for {
	}
}
