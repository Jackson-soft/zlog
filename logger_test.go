package zlog

import (
	"testing"
)

func TestNewZLog(t *testing.T) {
	b, err := NewInciseFile("xlog", "xlog.log", "xxlog", 500)
	if err != nil {
		t.Error(err)
	}
	SetBackend(b)
	Infoln("this is a message!!")
}

func BenchmarkLoops(t *testing.B) {
	z := NewZLog(InfoLevel)
	b, err := NewInciseFile("xlog", "xlog.log", "xxlog", 500)
	if err != nil {
		t.Error(err)
	}
	z.SetBackend(b)
	//z.AddBackend(b)
	for i := 0; i < t.N; i++ {
		z.output(WarnLevel, "sdfasdf\n ssfdsdfs \n asdfsd")
	}
}

func TestZLogLink(t *testing.T) {

}

func TestZLog(t *testing.T) {
	z := NewZLog(InfoLevel)
	b, err := NewInciseFile("xlog", "xlog.log", "xxlog", 500)
	if err != nil {
		t.Error(err)
	}
	z.SetBackend(b)
	//z.AddBackend(b)

	z.output(WarnLevel, "sdfasdf\n ssfdsdfs \n asdfsd")
}
