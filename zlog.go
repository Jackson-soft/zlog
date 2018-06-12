package zlog

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Fields type
type Fields map[string]interface{}

// LogLevel 日志等级
type LogLevel uint32

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Convert the Level to a string
func (level LogLevel) String() string {
	switch level {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	}

	return "unknown"
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (LogLevel, error) {
	switch strings.ToLower(lvl) {
	case "trace":
		return TraceLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l LogLevel
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

// ZLog is a log
type ZLog struct {
	mutex sync.Mutex

	formatter Formatter

	level  LogLevel
	buffer chan []byte

	stop chan bool

	backends []Backend
}

//NewZLog 创建日志
func NewZLog(level LogLevel) *ZLog {
	z := new(ZLog)

	z.formatter = new(TextFormatter)

	z.level = level
	z.stop = make(chan bool)
	z.buffer = make(chan []byte, 256)

	z.backends = []Backend{os.Stdout}

	go z.run()

	return z
}

// SetLevel 设置日志级别
func (z *ZLog) SetLevel(level LogLevel) {
	z.level = level
}

//SetFormattor 设置格式化前端
func (z *ZLog) SetFormattor(ft Formatter) {
	z.formatter = ft
}

//SetBackend 设置输出后端
func (z *ZLog) SetBackend(be Backend) {
	z.backends = []Backend{be}
}

//AddBackend 添加多个输出后端
func (z *ZLog) AddBackend(be Backend) {
	z.backends = append(z.backends, be)
}

//Stop 停止
func (z *ZLog) Stop() {
	z.stop <- true
}

func (z *ZLog) run() {
	for {
		select {
		case buf := <-z.buffer:
			for _, b := range z.backends {
				b.Write(buf)
			}
		case stop := <-z.stop:
			if stop && len(z.buffer) == 0 {
				for _, b := range z.backends {
					b.Close()
				}
				close(z.buffer)
				close(z.stop)
				return
			}
		}
	}
}

// Output 输出
func (z *ZLog) output(level LogLevel, msg string) {
	z.mutex.Lock()
	defer z.mutex.Unlock()
	if level >= z.level {
		buf := z.formatter.Format(level, msg)
		z.buffer <- buf
	}
}
