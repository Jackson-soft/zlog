package zlog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel uint32

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
	DayFormat  = "20060102"
)

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
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

// 日志格式 2006-01-02 15:04:05 info test.go 245 function : this is a error

// ZLog is a log
type ZLog struct {
	level       LogLevel
	out         io.Writer
	logPath     string // 文件存放目录
	maxFileSize int64  // 日志文件最大大小，单位M

	mutex       sync.Mutex
	logLocation string   // 文件的名称
	logIndex    int8     // 文件序号
	currentDay  string   // 当前日期
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

// SetLevel 设置日志级别
func (z *ZLog) SetLevel(level LogLevel) {
	z.level = level
}

// GetLevel 获取日志级别
func (z *ZLog) GetLevel() LogLevel {
	return z.level
}

//SetMaxFileSize 设置最大文件限制
func (z *ZLog) SetMaxFileSize(maxSize int64) {
	z.maxFileSize = maxSize
}

//SetLogPath 设置日志存放目录
func (z *ZLog) SetLogPath(logPath string) {
	z.logPath = logPath
}

// Output 输出
func (z *ZLog) Output(calldepth int, level LogLevel, msg string) error {
	if level >= z.level && z.checkFile() {
		pc, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		fn := runtime.FuncForPC(pc)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		logHead := fmt.Sprintf("[%s] %s %s %s %d", level.String(), time.Now().Format(TimeFormat), file, fn.Name(), line)
		buf := fmt.Sprintf("%s : %v \n", logHead, msg)

		_, err := z.out.Write([]byte(buf))
		return err
	}
	return nil
}

// createDir
func (z *ZLog) createDir(logPath string) error {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err := os.Mkdir(logPath, os.ModeDir|os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// openFile
func (z *ZLog) openFile() error {
	z.logLocation = fmt.Sprintf("%s/zlog-%s-%.4d.log", z.logPath, z.currentDay, z.logIndex)
	var err error
	z.out, err = os.OpenFile(z.logLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
		z.openFile()
	}
	return true
}

// getLogBlockSize 检查日志文件大小
func (z *ZLog) getLogBlockSize() (int64, error) {
	fileInfo, err := os.Stat(z.logLocation)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
