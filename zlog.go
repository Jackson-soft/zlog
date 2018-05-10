package zlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

//LogLevel 日志等级
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
	timeFormat = "2006-01-02 15:04:05"
	dayFormat  = "20060102"
)

var (
	defaultPath    = "zlog"
	defaultLink    = ""
	defaultMaxSize = int64(500 * 1024 * 1024) // 500M
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

// 日志格式 [error] 2006-01-02 15:04:05  test.go 245 function : this is a error

// ZLog is a log
type ZLog struct {
	level       LogLevel
	out         io.Writer
	logPath     string // 文件存放目录
	logLink     string // 文件软链接
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

	z.logPath = defaultPath
	z.logLink = defaultLink
	z.maxFileSize = defaultMaxSize
	z.logIndex = defaultIndex
	z.currentDay = time.Now().Format(dayFormat)
	z.logChang = true

	z.out = nil

	return z
}

// SetLevel 设置日志级别
func (z *ZLog) SetLevel(level LogLevel) {
	z.level = level
}

//SetMaxFileSize 设置最大文件限制,参数单位是M
func (z *ZLog) SetMaxFileSize(maxSize int64) {
	z.maxFileSize = maxSize * 1024 * 1024
}

//SetLogPath 设置日志存放目录
func (z *ZLog) SetLogPath(logPath string) {
	z.logPath = logPath
	z.logChang = true
}

//SetLogLink 设置最新日志文件的软链接
func (z *ZLog) SetLogLink(logLink string) {
	z.logLink = logLink
	z.logChang = true
}

// Output 输出
func (z *ZLog) Output(calldepth int, level LogLevel, msg string) error {
	z.mutex.Lock()
	defer z.mutex.Unlock()
	if level >= z.level {
		z.checkLogData()
		z.checkLogSize()
		if z.logChang {
			if err := z.openFile(); err != nil {
				return err
			}
		}
		pc, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		} else {
			slash := strings.LastIndex(file, "/")
			if slash >= 0 {
				file = file[slash+1:]
			}
		}
		fn := runtime.FuncForPC(pc)
		buf := fmt.Sprintf("%s [%s] %s %s %d : %s \n", time.Now().Format(timeFormat), level.String(), file, fn.Name(), line, msg)

		_, err := z.out.Write([]byte(buf))
		return err
	}
	return nil
}

// openFile 打开文件
func (z *ZLog) openFile() error {
	var err error
	if _, err = os.Stat(z.logPath); os.IsNotExist(err) {
		if err = os.Mkdir(z.logPath, os.ModeDir|os.ModePerm); err != nil {
			return err
		}
	}

	fileName := fmt.Sprintf("zlog-%s-%.4d.log", z.currentDay, z.logIndex)
	z.logLocation = filepath.Join(z.logPath, fileName)

	z.out, err = os.OpenFile(z.logLocation, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if z.logLink != "" {
		linkName := filepath.Join(z.logPath, z.logLink)
		tmpLinkName := linkName + `_symlink`
		if err := os.Symlink(fileName, tmpLinkName); err != nil {
			return err
		}

		if err := os.Rename(tmpLinkName, linkName); err != nil {
			return err
		}
	}
	z.logChang = false
	return nil
}

// checkLogData 检测日期是否与当前日期不一致
func (z *ZLog) checkLogData() {
	cDay := time.Now().Format(dayFormat)
	if cDay != z.currentDay {
		z.currentDay = cDay
		z.logChang = true
	}
}

// checkLogSize 用来检测日志文件是否超过规定大小
func (z *ZLog) checkLogSize() {
	fileInfo, err := os.Stat(z.logLocation)
	if err != nil {
		return
	}

	if fileInfo.Size() >= z.maxFileSize {
		z.logIndex++
		z.logChang = true
	}
}
