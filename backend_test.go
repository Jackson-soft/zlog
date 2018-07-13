package zlog

import (
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	bb, err := NewInciseFile("xlog", "", "", 0)
	if err != nil {
		b.Error(err)
	}
	msg := []byte("this is a message!!!\n")
	for i := 0; i < b.N; i++ {
		bb.Write(msg)
	}
}

func TestWrite(t *testing.T) {
	bb, err := NewInciseFile("xlog", "", "", 0)
	if err != nil {
		t.Error(err)
	}
	msg := []byte("this is a message!!!\n")

	bb.Write(msg)
}
