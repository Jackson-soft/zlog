package zlog

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// LogLevel 日志等级
type Level uint8

const (
	TraceLevel Level = iota
	DebugLevel
	InforLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	NULLLevel //非法等级
)

// Convert the Level to a string
func (level Level) String() string {
	switch level {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InforLevel:
		return "infor"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return ""
	}
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "trace":
		return TraceLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn":
		return WarnLevel, nil
	case "infor":
		return InforLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	return NULLLevel, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

// ZLog is a log
type ZLog struct {
	mutex sync.Mutex

	formatter Formatter
	backends  []Backend
	buffer    chan []byte

	stop chan bool
}

//NewZLog 创建日志
func NewZLog(level Level) *ZLog {
	z := new(ZLog)

	z.formatter = &TextFormatter{
		strLevel: level.String(),
		tLevel:   level,
		data:     make(Fields),
	}

	z.stop = make(chan bool)
	z.buffer = make(chan []byte, 256)

	z.backends = []Backend{os.Stdout}

	go z.run()

	return z
}

// SetLevel 设置日志级别
func (z *ZLog) SetLevel(level string) error {
	return z.formatter.SetLevel(level)
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
func (z *ZLog) output(level Level, msg string) {
	z.mutex.Lock()
	defer z.mutex.Unlock()
	buf := z.formatter.Format(level, msg)
	if len(buf) > 0 {
		z.buffer <- buf
	}
}

func (z *ZLog) WithFields(fields Fields) *ZLog {
	z.formatter.WithFields(fields)
	return z
}

// Tracef logs a message at level Info on the standard logger.
func (z *ZLog) Tracef(format string, args ...interface{}) {
	z.output(TraceLevel, fmt.Sprintf(format, args...))
}

// Debugf logs a message at level Debug on the standard logger.
func (z *ZLog) Debugf(format string, args ...interface{}) {
	z.output(DebugLevel, fmt.Sprintf(format, args...))
}

// Infof logs a message at level Info on the standard logger.
func (z *ZLog) Infof(format string, args ...interface{}) {
	z.output(InforLevel, fmt.Sprintf(format, args...))
}

// Warnf logs a message at level Warn on the standard logger.
func (z *ZLog) Warnf(format string, args ...interface{}) {
	z.output(WarnLevel, fmt.Sprintf(format, args...))
}

// Errorf logs a message at level Error on the standard logger.
func (z *ZLog) Errorf(format string, args ...interface{}) {
	z.output(ErrorLevel, fmt.Sprintf(format, args...))
}

// Fatalf logs a message at level Fatal on the standard logger.
func (z *ZLog) Fatalf(format string, args ...interface{}) {
	z.output(FatalLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Debugln logs a message at level Debug on the standard logger.
func (z *ZLog) Debugln(args ...interface{}) {
	z.output(DebugLevel, fmt.Sprint(args...))
}

// Infoln logs a message at level Info on the standard logger.
func (z *ZLog) Infoln(args ...interface{}) {
	z.output(InforLevel, fmt.Sprint(args...))
}

// Warnln logs a message at level Warn on the standard logger.
func (z *ZLog) Warnln(args ...interface{}) {
	z.output(WarnLevel, fmt.Sprint(args...))
}

// Errorln logs a message at level Error on the standard logger.
func (z *ZLog) Errorln(args ...interface{}) {
	z.output(ErrorLevel, fmt.Sprint(args...))
}

// Fatalln logs a message at level Fatal on the standard logger.
func (z *ZLog) Fatalln(args ...interface{}) {
	z.output(FatalLevel, fmt.Sprint(args...))
	os.Exit(1)
}
