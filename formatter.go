package zlog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Fields type
type Fields map[string]interface{}

//Formatter 格式化前端
type Formatter interface {
	Format(level Level, msg string) []byte
	WithFields(fields Fields)
}

//TextFormatter 文本格式化前端
type TextFormatter struct {
	Data Fields
}

//Format 日志格式  2006-01-02 15:04:05 [error] test.go 245 : this is a error
func (f *TextFormatter) Format(level Level, msg string) []byte {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	} else {
		slash := strings.LastIndex(file, string(filepath.Separator))
		if slash >= 0 {
			file = file[slash+1:]
		}
	}

	buf := &strings.Builder{}

	fmt.Fprintf(buf, "%s [%s] %s %d :: ", time.Now().Format(timeFormat), level.String(), file, line)

	if len(f.Data) > 0 {
		for k, v := range f.Data {
			fmt.Fprintf(buf, "%v:%v,", k, v)
		}
	}

	fmt.Fprintf(buf, "%s\n", msg)

	return []byte(buf.String())
}

//WithFields 添加一个map
func (f *TextFormatter) WithFields(fields Fields) {
	if len(fields) == 0 {
		return
	}
	for k, v := range fields {
		f.Data[k] = v
	}
}
