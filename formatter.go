package zlog

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

//Formatter 格式化前端
type Formatter interface {
	Format(level LogLevel, msg string) []byte
}

//TextFormatter 文本格式化前端
type TextFormatter struct{}

//Format 日志格式  2006-01-02 15:04:05 [error] test.go 245 function : this is a error
func (f *TextFormatter) Format(level LogLevel, msg string) []byte {
	pc, file, line, ok := runtime.Caller(3)
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
	return []byte(buf)
}
