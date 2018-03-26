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

var (
	defaultLogPath = "zlog"
	defaultMaxSize = int64(500) // 500M
	defaultIndex   = 1
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
	logLocation string // 文件的名称
	logIndex    int    // 文件序号
	currentDay  string // 当前日期
	logChang    bool   //日志文件是否要切割
}

//NewZLog 创建日志
func NewZLog(level LogLevel) *ZLog {
	z := new(ZLog)

	z.level = level

	z.logPath = defaultLogPath
	z.maxFileSize = defaultMaxSize
	z.logIndex = defaultIndex
	z.currentDay = time.Now().Format(DayFormat)
	z.logChang = true

	z.out = nil

	return z
}

// SetLevel 设置日志级别
func (z *ZLog) SetLevel(level LogLevel) {
	z.level = level
}

//SetMaxFileSize 设置最大文件限制
func (z *ZLog) SetMaxFileSize(maxSize int64) {
	z.maxFileSize = maxSize
}

//SetLogPath 设置日志存放目录
func (z *ZLog) SetLogPath(logPath string) {
	z.logPath = logPath
	z.logChang = true
}

// Output 输出
func (z *ZLog) Output(calldepth int, level LogLevel, msg string) error {
	z.mutex.Lock()
	defer z.mutex.Unlock()
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

// openFile
func (z *ZLog) openFile() error {
	if _, err := os.Stat(z.logPath); os.IsNotExist(err) {
		if err := os.Mkdir(z.logPath, os.ModeDir|os.ModePerm); err != nil {
			return err
		}
	}

	z.logLocation = fmt.Sprintf("%s/zlog-%s-%.4d.log", z.logPath, z.currentDay, z.logIndex)
	var err error
	z.out, err = os.OpenFile(z.logLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	z.logChang = false
	return nil
}

// checkFile 用来检测日志文件是否超过规定大小
func (z *ZLog) checkFile() bool {
	fileSize, err := z.getLogBlockSize()
	if err != nil {
		return false
	}

	bChang := false
	// 文件超过大小或日期不是同一天
	if fileSize >= z.maxFileSize*1024*1024 {
		z.logIndex++
		bChang = true
	}

	cDay := time.Now().Format(DayFormat)
	if cDay != z.currentDay {
		z.currentDay = cDay
		bChang = true
	}
	if bChang {
		if err = z.openFile(); err != nil {
			return false
		}
	}
	return true
}

// getLogBlockSize 检查日志文件大小
func (z *ZLog) getLogBlockSize() (int64, error) {
	if z.logChang {
		if err := z.openFile(); err != nil {
			return 0, err
		}
	}
	fileInfo, err := os.Stat(z.logLocation)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
