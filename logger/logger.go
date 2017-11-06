// Package logger implements a very basic logger
//
// currently only logging to stdout are implemented
// NEXT: maybe log all messages in an extra chat channel 'matterbot-logs'?
//
//
// ------------------------------------------------------------------------
// why this package and not the 'log' package from the stdlib?
//
// to use log-level prefixes like 'DEBUG, INFO, ...' and enable / disable different
// levels, you need to instance an logger for each level 'with log.New(out, prefix, flag)'.
// if you make this in the main package, the code from the 'chat' and 'mail'
// packages can't use this instances.
package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

// to prevent interleaved log output
var mutex = sync.Mutex{}

func log(level LogLevel, xs ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	if level < logLevel {
		return
	}

	pos := "???"
	if _, file, line, ok := runtime.Caller(2); ok {
		fileName := path.Base(file)
		fileNameWithoutExt := strings.TrimRight(fileName, path.Ext(fileName))
		pos = fmt.Sprintf("%10s:%3d", fileNameWithoutExt, line)
	}

	fmt.Printf("%s %5s %s | ", time.Now().Format("02.01 15:04:05.000"), level.String(), pos)
	fmt.Println(xs...)
}

func Debug(xs ...interface{}) {
	log(DebugLevel, xs...)
}
func Debugf(format string, xs ...interface{}) {
	log(DebugLevel, fmt.Sprintf(format, xs...))
}
func Info(xs ...interface{}) {
	log(InfoLevel, xs...)
}
func Infof(format string, xs ...interface{}) {
	log(InfoLevel, fmt.Sprintf(format, xs...))
}
func Error(xs ...interface{}) {
	log(ErrorLevel, xs...)
}
func Errorf(format string, xs ...interface{}) {
	log(ErrorLevel, fmt.Sprintf(format, xs...))
}
