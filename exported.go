package zlog

var std = NewZLog(InforLevel)

func GetInstance() *ZLog {
	return std
}

// SetLevel 设置日志等级
func SetLevel(lvl string) error {
	std.mutex.Lock()
	err := std.SetLevel(lvl)
	std.mutex.Unlock()
	return err
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

func Stop() {
	std.mutex.Lock()
	std.Stop()
	std.mutex.Unlock()
}

func WithFields(fields Fields) *ZLog {
	return std.WithFields(fields)
}

// Tracef logs a message at level Info on the standard logger.
func Tracef(format string, args ...interface{}) {
	std.Tracef(format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	std.Debugln(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	std.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	std.Warnln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	std.Errorln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	std.Fatalln(args...)
}
