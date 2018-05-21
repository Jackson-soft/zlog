package zlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

//Backend 日志输出端
type Backend interface {
	Write(buf []byte) (int, error)
	Close() error
}

//InciseFileBackend 文件切割后端
type InciseFileBackend struct {
	fd *os.File

	filePath    string // 文件存放目录
	fileLink    string // 文件软链接
	namePrefix  string // 日志文件名前缀
	maxFileSize int64  // 日志文件最大大小，单位M

	appellation string // 文件的名称
	index       int    // 文件序号
	currentDay  string // 当前日期
	chang       bool   // 日志文件是否要切割
}

const (
	timeFormat = "2006-01-02 15:04:05"
	dayFormat  = "20060102"

	defaultPath    = ""
	defaultLink    = ""
	defaultPrefix  = ""
	defaultMaxSize = int64(500 * 1024 * 1024) // 500M
	defaultIndex   = 1
)

func NewInciseFile(filePath, fileLink, prefix string, maxSize int64) (*InciseFileBackend, error) {
	b := new(InciseFileBackend)
	var err error
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		if err = os.Mkdir(filePath, os.ModeDir|os.ModePerm); err != nil {
			return nil, err
		}
	}
	b.filePath = filePath

	b.fileLink = fileLink
	b.namePrefix = prefix
	b.maxFileSize = maxSize * 1014 * 1024

	b.currentDay = time.Now().Format(dayFormat)
	b.index = defaultIndex

	b.chang = true

	return b, nil
}

func (b *InciseFileBackend) doIncise() error {
	b.checkData()
	b.checkSize()
	if b.chang {
		fileName := fmt.Sprintf("%s-%s-%.4d.log", b.namePrefix, b.currentDay, b.index)
		b.appellation = filepath.Join(b.filePath, fileName)
		var err error
		b.fd, err = os.OpenFile(b.appellation, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		if b.fileLink != "" {
			linkName := filepath.Join(b.filePath, b.fileLink)
			tmpLinkName := linkName + `_symlink`
			if err = os.Symlink(fileName, tmpLinkName); err != nil {
				return err
			}

			if err = os.Rename(tmpLinkName, linkName); err != nil {
				return err
			}
		}
		b.chang = false
	}
	return nil
}

func (b *InciseFileBackend) checkData() {
	cDay := time.Now().Format(dayFormat)
	if cDay != b.currentDay {
		b.currentDay = cDay
		b.chang = true
	}
}

func (b *InciseFileBackend) checkSize() {
	fileInfo, err := os.Stat(b.appellation)
	if err != nil {
		return
	}

	if fileInfo.Size() >= b.maxFileSize {
		b.index++
		b.chang = true
	}
}

//Write 文件后端写方法
func (b *InciseFileBackend) Write(buf []byte) (int, error) {
	b.doIncise()
	return b.fd.Write(buf)
}

//Close 文件后端关闭
func (b *InciseFileBackend) Close() error {
	return b.fd.Close()
}
