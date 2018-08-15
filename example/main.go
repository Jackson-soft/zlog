package main

import (
	"runtime/debug"
	"zlog"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Errorf("%s\n", debug.Stack())
			zlog.Stop()
		}
	}()
	zlog.Infoln("df")
	zlog.WithFields(zlog.Fields{"fff": "dff", "vvv": 45.6}).Infoln("dfsdfa")

	panic("dfdf")
	for {
	}
}
