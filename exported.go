package zlog

import "os"

var std = NewZLog()

func GetInstance() *ZLog {
	if std == nil {
		return NewZLog()
	}
	return std
}

// SetLevel sets the standard logger level.
func SetLevel(level LogLevel) {
	std.mutex.Lock()
	std.SetLevel(level)
	std.mutex.Unlock()
}

// GetLevel returns the standard logger level.
func GetLevel() LogLevel {
	std.mutex.Lock()
	defer std.mutex.Unlock()
	return std.GetLevel()
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
