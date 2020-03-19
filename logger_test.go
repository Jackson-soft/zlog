package zlog

import (
	"testing"
)

func TestLevel(t *testing.T) {
	lvl, err := ParseLevel("infor")
	if err != nil {
		t.Error(err)
	}
	t.Log(lvl.String())
}

func TestNewZLog(t *testing.T) {
	b, err := NewInciseFile("xlog", "xlog.log", "xxlog", 500)
	if err != nil {
		t.Error(err)
	}
	SetBackend(b)
	Infoln("this is a message!!")
}

func BenchmarkLoops(t *testing.B) {
	z := NewZLog(InforLevel)
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

func TestZLog(t *testing.T) {
	z := NewZLog(InforLevel)
	b, err := NewInciseFile("xlog", "xlog.log", "xxlog", 500)
	if err != nil {
		t.Error(err)
	}
	z.SetBackend(b)
	//z.AddBackend(b)

	z.output(WarnLevel, "sdfasdf\n ssfdsdfs \n asdfsd")
}
