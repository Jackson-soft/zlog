package zlog

import "os"

var std = NewZLog(InfoLevel)

func GetInstance() *ZLog {
	return std
}

// SetLevel sets the standard logger level.
func SetLevel(lvl string) error {
	std.mutex.Lock()
	level, err := ParseLevel(lvl)
	if err != nil {
		return err
	}
	std.SetLevel(level)
	std.mutex.Unlock()
	return nil
}

//SetMaxFileSize 设置最大文件限制
func SetMaxFileSize(maxSize int64) {
	std.mutex.Lock()
	std.SetMaxFileSize(maxSize)
	std.mutex.Unlock()
}

//SetLogPath 设置日志存放目录
func SetLogPath(logPath string) {
	std.mutex.Lock()
	std.SetLogPath(logPath)
	std.mutex.Unlock()
}

//SetLogLink 设置日志软链接
func SetLogLink(logLink string) {
	std.mutex.Lock()
	std.SetLogLink(logLink)
	std.mutex.Unlock()
}

func LogTrace(msg string) {
	std.Output(2, TraceLevel, msg)
}

func LogDebug(msg string) {
	std.Output(2, DebugLevel, msg)
}

func LogInfo(msg string) {
	std.Output(2, InfoLevel, msg)
}

func LogWarn(msg string) {
	std.Output(2, WarnLevel, msg)
}

func LogError(msg string) {
	std.Output(2, ErrorLevel, msg)
}

func LogFatal(msg string) {
	std.Output(2, FatalLevel, msg)
	os.Exit(1)
}
