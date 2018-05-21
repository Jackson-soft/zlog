package zlog

import (
	"testing"
)

func TestNewZLog(t *testing.T) {
}

func BenchmarkLoops(t *testing.B) {
	z := NewZLog(InfoLevel)
	b, err := NewInciseFile("xlog", "xlog.log", "xxlog", 500)
	if err != nil {
		t.Error(err)
	}
	//SetBackend(b)
	z.AddBackend(b)
	for i := 0; i < t.N; i++ {
		//Warnln("sdfasdf\n ssfdsdfs \n asdfsd")
		z.output(WarnLevel, "sdfasdf\n ssfdsdfs \n asdfsd")
	}
}

func TestZLogLink(t *testing.T) {

}

func TestZLog(t *testing.T) {

}
