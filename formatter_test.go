package zlog

import "testing"

func BenchmarkFormattor(b *testing.B) {
	tf := TextFormatter{}
	for i := 0; i < b.N; i++ {
		tf.Format(InforLevel, "this is a message!!!")
	}
}
