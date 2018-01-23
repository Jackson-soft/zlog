package zlog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel uint8

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	DayFormat  = "2006-01-02"
)

// 日志格式 2006-01-02 15:04:05 info test.go 245 function : this is a error

// ZLog is a log
type ZLog struct {
	mutex       sync.Mutex
	level       LogLevel
	logPath     string // 文件存放目录
	maxFileSize int64  // 日志文件最大大小，单位M
	logLocation string // 文件的名称
	logIndex    int8   // 文件序号
	logfd       *os.File
	currentDay  string   // 当前时期
	buffer      [][]byte // 这里需要一个环形缓冲区
}

func NewZLog() *ZLog {
	l := new(ZLog)
	l.level = InfoLevel
	l.logPath = "zlog"
	l.maxFileSize = 500
	l.logIndex = 1
	l.currentDay = time.Now().Format(DayFormat)
	l.buffer = make([][]byte, 0)
	if err := l.createDir(l.logPath); err != nil {
		return nil
	}
	if err := l.openFile(); err != nil {
		return nil
	}
	return l
}

func (z *ZLog) Run() {
	for {
		if len(z.buffer) > 0 && z.checkFile() {
			z.logfd.Write(z.popBuffer())
		}
	}
}

// Initialization 日志的初始化
func (z *ZLog) Initialization() {

}

// SetLevel 设置日志级别
func (z *ZLog) SetLevel(level LogLevel) {
	z.level = level
}

// GetLevel 获取日志级别
func (z *ZLog) GetLevel() LogLevel {
	return z.level
}

func (z *ZLog) Close() {
	z.logfd.Close()
}

// logLevelToString
func (z *ZLog) logLevelToString(level LogLevel) string {
	switch level {
	case TraceLevel:
		return "Trace"
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case WarnLevel:
		return "Warn"
	case ErrorLevel:
		return "Error"
	case FatalLevel:
		return "Fatal"
	default:
		return "Debug"
	}
}

// stringToLogLevel
func (z *ZLog) stringToLogLevel(str string) LogLevel {
	level := strings.ToLower(str)
	switch level {
	case "trace":
		return TraceLevel
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return DebugLevel
	}
}

// createDir
func (z *ZLog) createDir(logPath string) error {
	return os.Mkdir(logPath, os.ModeDir)
}

// openFile
func (z *ZLog) openFile() error {
	z.logLocation = fmt.Sprintf("%s/zlog-%s-%d.log", z.logPath, z.currentDay, z.logIndex)
	var err error
	z.logfd, err = os.OpenFile(z.logLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	return err
}

// checkFile 用来检测日志文件是否超过规定大小
func (z *ZLog) checkFile() bool {
	fileSize, err := z.getLogBlockSize()
	if err != nil {
		return false
	}
	cDay := time.Now().Format(DayFormat)
	// 文件超过大小或日期不是同一天
	if fileSize >= z.maxFileSize*1024*1024 || cDay != z.currentDay {
		z.logIndex++
		z.currentDay = cDay
		z.logfd.Close()
		z.openFile()
	}
	return true
}

// prepareLogHead 组装日志头
func (z *ZLog) prepareLogHead(level LogLevel) string {
	pc, file, line, ok := runtime.Caller(0)
	if !ok {
		file = "???"
		line = 0
	}
	fn := runtime.FuncForPC(pc)

	return fmt.Sprintf("%s %s %s %s %d", time.Now().Format(TimeFormat), z.logLevelToString(level), file, fn.Name(), line)
}

// getLogBlockSize 检查日志文件大小
func (z *ZLog) getLogBlockSize() (int64, error) {
	fileInfo, err := os.Stat(z.logLocation)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// pushBuffer 压入
func (z *ZLog) pushBuffer(buf []byte) {
	//z.mutex.Lock()
	z.buffer = append(z.buffer, buf)
	//z.mutex.Unlock()
}

// popBuffer 弹出
func (z *ZLog) popBuffer() []byte {
	//z.mutex.Lock()

	buf := new([]byte)
	if len(z.buffer) == 0 {
		*buf = make([]byte, 0)
	} else {
		*buf = z.buffer[0]
		z.buffer = z.buffer[1:]
	}

	//z.mutex.Unlock()
	return *buf
}

// Output 输出
func (z *ZLog) Output(level LogLevel, msg string) {
	if level >= z.level {
		logHead := z.prepareLogHead(level)
		buf := fmt.Sprintf("%s : %s", logHead, msg)
		z.pushBuffer([]byte(buf))
	}
}
