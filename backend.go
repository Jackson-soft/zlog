package zlog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

//Backend 日志输出端接口
type Backend interface {
	Write(buf []byte) (int, error)
	Sync() error
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
	changed     bool   // 日志文件是否要切割
}

const (
	timeFormat = "2006-01-02 15:04:05.000"
	dayFormat  = "20060102"

	defaultMaxSize = int64(500 * 1024 * 1024) // 500M
	defaultIndex   = 1

	defaultPrefix = "zlog"
	defaultSuffix = ".log"
)

//NewInciseFile create file backend
func NewInciseFile(filePath, fileLink, prefix string, maxSize int64) (*InciseFileBackend, error) {
	if len(filePath) == 0 {
		return nil, errors.New("file path is nil")
	}

	var err error
	if _, err = os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir(filePath, os.ModeDir|os.ModePerm); err != nil {
			return nil, err
		}
	}

	b := new(InciseFileBackend)

	b.filePath = filePath
	b.currentDay = time.Now().Format(dayFormat)
	b.index = defaultIndex
	b.fileLink = fileLink
	b.changed = true

	if len(prefix) == 0 {
		b.namePrefix = defaultPrefix
	} else {
		b.namePrefix = prefix
	}

	if maxSize == 0 {
		b.maxFileSize = defaultMaxSize
	} else {
		b.maxFileSize = maxSize * 1014 * 1024
	}

	if err = b.createFile(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *InciseFileBackend) doIncise() error {
	b.checkData()
	b.checkSize()
	if b.changed {
		return b.createFile()
	}
	return nil
}

func (b *InciseFileBackend) checkData() {
	cDay := time.Now().Format(dayFormat)
	if cDay != b.currentDay {
		b.currentDay = cDay
		//日期变更后，序号重置
		b.index = defaultIndex
		b.changed = true
	}
}

func (b *InciseFileBackend) checkSize() {
	fileInfo, err := os.Stat(b.appellation)
	if err != nil {
		return
	}

	if fileInfo.Size() >= b.maxFileSize {
		b.index++
		b.changed = true
	}
}

//Write 文件后端写方法
func (b *InciseFileBackend) Write(buf []byte) (int, error) {
	if err := b.doIncise(); err != nil {
		return 0, err
	}
	return b.fd.Write(buf)
}

//Sync 刷新缓存区
func (b *InciseFileBackend) Sync() error {
	return b.fd.Sync()
}

//Close 文件后端关闭
func (b *InciseFileBackend) Close() error {
	return b.fd.Close()
}

//createFile 格式化文件名字
func (b *InciseFileBackend) createFile() error {
	fileName := fmt.Sprintf("%s-%s-%.4d%s", b.namePrefix, b.currentDay, b.index, defaultSuffix)
	appellation := filepath.Join(b.filePath, fileName)

	fd, err := os.OpenFile(appellation, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 新打开文件之前先关闭，不然会有文件描述符泄漏
	if b.fd != nil {
		// flush
		b.fd.Sync()
		b.fd.Close()
		b.fd = nil
	}

	b.fd = fd
	b.appellation = appellation

	//如果软链接配置不为空
	if len(b.fileLink) > 0 {
		linkName := filepath.Join(b.filePath, b.fileLink)
		tmpLinkName := linkName + `_symlink`
		if err = os.Symlink(fileName, tmpLinkName); err != nil {
			return err
		}

		if err = os.Rename(tmpLinkName, linkName); err != nil {
			return err
		}
	}
	b.changed = false
	return nil
}
