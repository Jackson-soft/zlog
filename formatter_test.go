package zlog

import "testing"

func BenchmarkFormattor(b *testing.B) {
	tf := TextFormatter{}
	for i := 0; i < b.N; i++ {
		tf.Format(InfoLevel, "this is a message!!!")
	}
}
