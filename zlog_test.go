package zlog

import (
	"testing"
)

func TestNewZLog(t *testing.T) {
	z := NewZLog(InfoLevel)
	z.Output(InfoLevel, "dfsfa")
}

func TestZLogLink(t *testing.T) {

}

func TestZLog(t *testing.T) {

}
