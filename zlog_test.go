package zlog

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestNewZLog(t *testing.T) {
	z := NewZLog(InfoLevel)
	z.SetLogPath("xlog")
	z.Output(2, ErrorLevel, "dfasdf")
}

func TestZLogLink(t *testing.T) {
	z := NewZLog(InfoLevel)
	z.SetLogLink("xlog.log")
	z.Output(2, ErrorLevel, "fdsfasd")
	d, err := ioutil.ReadFile("xlog.log")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(d))
}

func TestZLog(t *testing.T) {
	SetLevel("info")
	SetLogLink("xx.log")
	LogError("dfasdf")
}
