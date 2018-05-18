package zlog

import (
	"fmt"
	"os"
)

var std = NewZLog(InfoLevel)

func GetInstance() *ZLog {
	return std
}

// SetLevel 设置日志等级
func SetLevel(lvl string) error {
	std.mutex.Lock()
	level, err := ParseLevel(lvl)
	if err != nil {
		return err
	}
	std.level = level
	std.mutex.Unlock()
	return nil
}

//SetFormattor 设置格式化器
func SetFormattor(ft Formatter) {
	std.mutex.Lock()
	std.formatter = ft
	std.mutex.Unlock()
}

//AddBackend 添加多个输出后端
func AddBackend(be Backend) {
	std.mutex.Lock()
	std.backends = append(std.backends, be)
	std.mutex.Unlock()
}

//SetBackend 设置输出后端为单一
func SetBackend(be Backend) {
	std.mutex.Lock()
	std.backends = []Backend{be}
	std.mutex.Unlock()
}

// Tracef logs a message at level Info on the standard logger.
func Tracef(format string, args ...interface{}) {
	std.output(TraceLevel, fmt.Sprintf(format, args...))
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	std.output(DebugLevel, fmt.Sprintf(format, args...))
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	std.output(InfoLevel, fmt.Sprintf(format, args...))
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	std.output(WarnLevel, fmt.Sprintf(format, args...))
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	std.output(ErrorLevel, fmt.Sprintf(format, args...))
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	std.output(FatalLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	std.output(DebugLevel, fmt.Sprint(args...))
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	std.output(InfoLevel, fmt.Sprint(args...))
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	std.output(WarnLevel, fmt.Sprint(args...))
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	std.output(ErrorLevel, fmt.Sprint(args...))
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	std.output(FatalLevel, fmt.Sprint(args...))
	os.Exit(1)
}
