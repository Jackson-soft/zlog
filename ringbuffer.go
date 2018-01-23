package zlog

import "sync"

const minLen = 16

type node struct {
}

type RingBuffer struct {
	mutex              sync.Mutex
	data               []interface{}
	write, read, count int
}

// New 默认环形缓存是16
func New(len int) *RingBuffer {
	r := new(RingBuffer)

	r.write = 0
	r.read = 0
	if len == 0 {
		len = minLen
	}
	r.count = len
	r.data = make([]interface{}, r.count)

	return r
}

//Size 缓存目前的大小
func (r *RingBuffer) Size() int {
	return r.count
}

//Empty 是否为空
func (r *RingBuffer) Empty() bool {
	return r.read == r.write
}

//resize 重新分配大小
func (r *RingBuffer) resize() {

}

// Push 压入
func (r *RingBuffer) Push(value interface{}) {
	r.mutex.Lock()
	r.write++
	r.data = append(r.data, value)
	r.mutex.Unlock()
}

//Pop 弹出
func (r *RingBuffer) Pop() interface{} {
	r.mutex.Lock()

	r.read++
	r.mutex.Unlock()
}
