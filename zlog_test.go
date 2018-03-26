package zlog

import "testing"

func TestNewZLog(t *testing.T) {
	z := NewZLog(InfoLevel)
	z.SetLogPath("xlog")
	z.Output(2, ErrorLevel, "dfasdf")
}
