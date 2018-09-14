package main

import (
	"runtime/debug"
	"time"

	"github.com/Jackson-soft/zlog"
)

func main() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zlog.Errorf("%s\n", debug.Stack())
				zlog.Sync()
			}
		}()
		panic("dfdf")
	}()
	time.Sleep(1 * time.Second)
	zlog.Infoln("df")
	zlog.WithFields(zlog.Fields{"fff": "dff", "vvv": 45.6}).Infoln("dfsdfa")

	for {
	}
}
