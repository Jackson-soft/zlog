package zlog

var std = NewZLog()

func GetInstance() *ZLog {
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
	std.Output(TraceLevel, msg)
}
func LogDebug(msg string) {
	std.Output(DebugLevel, msg)
}
func LogInfo(msg string) {
	std.Output(InfoLevel, msg)
}
func LogWarn(msg string) {
	std.Output(WarnLevel, msg)
}
func LogError(msg string) {
	std.Output(ErrorLevel, msg)
}
func LogFatal(msg string) {
	std.Output(FatalLevel, msg)
}
