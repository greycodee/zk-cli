package log

import (
	"bytes"
	"fmt"
	"time"
)

type Log struct {
	Logs bytes.Buffer
}

func NewLogs() *Log {
	return &Log{Logs: bytes.Buffer{}}
}

func (l *Log) Info(s string)  {
	l.Logs.WriteString(fmt.Sprintf("%s：%s\n",time.Stamp,s))
}