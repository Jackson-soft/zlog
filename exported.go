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
	os.Exit(1)
}
