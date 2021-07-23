package epilog

import (
	"fmt"
	"runtime"
	"time"
)

func innerPrint(lgr *Logger, newLine bool, level logLevel, args ...interface{}) {
	var content string
	if newLine {
		content = fmt.Sprintln(args...)
	} else {
		content = fmt.Sprint(args...)
	}
	_, filename, line, _ := runtime.Caller(2)

	var typ string
	switch level {
	case NORMAL:
		typ = "NORMAL"
	case DEBUG:
		typ = "DEBUG"
	case INFO:
		typ = "INFO"
	case WARNING:
		typ = "WARNING"
	case FATAL:
		typ = "FATAL"
	default:
		typ = "NORMAL"
	}
	lgr.buffer.Put(BufferItem{
		id: lgr.seqNum, typ: typ,
		filename: filename, line: line,
		time: time.Now(), content: content,
	})
}
