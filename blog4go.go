// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	"fmt"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

var (
	levelStrings = [...]string{"DEBUG", "TRAC", "INFO", "WARN", "ERROR", "CRITAL"}
)

func (self Level) ToString() string {
	if 0 > self || int(self) >= len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[self]
}

// 单条日志记录结构体
type LogRecord struct {
	level   Level
	message string
}

// 各种日志结构接口
type LogWriter interface {
	// 用户调用开始装逼地写log
	Start()

	// 提供用户主动将log输出到文件的方法
	// 当chan为有缓冲时
	Flush()

	// 关闭log writer的处理方法
	// 善后
	Close()

	// 用于内部写log的方法
	write(record *LogRecord)

	// 装逼logger自动将log输出到文件
	run()
}

// DefaultFileLogWriter.c 为无缓冲channel
var DefaultFileLogWriter *FileLogWriter = new(FileLogWriter)

var (
	DefaultBufferSize = 32
)

// 装逼的logger
type FileLogWriter struct {
	level Level

	c chan *LogRecord
	// c channel buffer size
	bufferSize int

	// log文件
	filename string
	file     *os.File

	// logrotate
	rotate bool
}

// 包初始化函数
func init() {

}

func (self *FileLogWriter) validateConfig() {
	if self.level < DEBUG || self.level > CRITICAL {
		panic("Please set an valid log level.")
	}

	if self.bufferSize < 0 {
		self.bufferSize = DefaultBufferSize
	}
}

func (self *FileLogWriter) Start() {
	self.validateConfig()

	self.c = make(chan *LogRecord, self.bufferSize)

	// 打开文件描述符

	go self.run()
}

func (self *FileLogWriter) write(record *LogRecord) {
	self.c <- record
}

func (self *FileLogWriter) Flush() {

}

func (self *FileLogWriter) Close() {
	close(self.c)
}

func (self *FileLogWriter) run() {
	for {
		select {
		case record := <-self.c:
			os.Stdout.WriteString(record.message)
		}
	}

}

func (self *FileLogWriter) SetLevel(level Level) *FileLogWriter {
	self.level = level
	return self
}

func (self *FileLogWriter) GetLevel() Level {
	return self.level
}

func (self *FileLogWriter) Debug(message string) {

}

func (self *FileLogWriter) Debugf(message string, args ...interface{}) {
	// 格式化构造message
	// 使用 % 作占位符

	// 识别占位符标记
	var tag bool = false

	for i := 0; i < len(message); i++ {
		switch message[i] {
		//占位符，百分号
		case '%':
			tag = true
		//占位符，有意义部分
		case 's':
			if !tag {
				continue
			}
		//转义符
		case '\\':
		}
	}
	fmt.Println(tag)
}
