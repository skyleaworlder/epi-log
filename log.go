package epilog

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	levelDefault      logLevel = INFO
	maxItemNumDefault int      = 1
)

// Logger is a struct
type Logger struct {
	// FATAL > WARNING > INFO > DEBUG
	level logLevel

	// map is much better than slice!
	// do not try to persuade me.
	appenders map[string]Appender
	buffer    *Buffer

	// seqNum(sequence number) is logItem id
	seqNum int

	mtx *sync.Mutex
}

// New is a constructor
func New() (lgr *Logger) {
	appenders := make(map[string]Appender)
	appenders["stdio"] = &StdAppender{}
	return &Logger{
		level:     levelDefault,
		appenders: appenders,
		buffer:    NewBuffer(maxItemNumDefault),
		seqNum:    0,
		mtx:       &sync.Mutex{},
	}
}

// Use is to use epilog.
// and Use method will generate a signal procecss goroutine.
// by this signal-process goroutine, epilog can process unexpected signals.
// what's more, Use also calls Monitor!
func (lgr *Logger) Use() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGABRT, syscall.SIGALRM, syscall.SIGBUS, syscall.SIGFPE,
		syscall.SIGHUP, syscall.SIGILL, syscall.SIGINT, syscall.SIGKILL,
		syscall.SIGPIPE, syscall.SIGQUIT, syscall.SIGSEGV, syscall.SIGTERM,
		syscall.SIGTRAP,
	)
	go func() {
		for sig := range sigs {
			switch sig {
			case syscall.SIGINT:
				lgr.Warningln("syscall.SIGINT:", sig)
			case syscall.SIGTERM:
				lgr.Warningln("syscall.SIGTERM:", sig)
			case syscall.SIGKILL:
				lgr.Warningln("syscall.SIGKILL:", sig)
			default:
				lgr.Warningln("I do not know what happened just now.")
			}
			lgr.End()
		}
	}()

	lgr.Monitor()
}

// End is to exit
func (lgr *Logger) End() {
	for item := range lgr.buffer.items {
		// append content to appenders
		for _, appender := range lgr.appenders {
			err := appender.Append(item.Serialize())
			if err != nil {
				// in End, "Append Failure" is a fatal error
				os.Exit(0)
				goto AppendError
			}
		}
		if len(lgr.buffer.items) == 0 {
			// should not use "lgr.buffer.empty <- true" here.
			// this statement is used in Monitor method.
			// so "lgr.buffer.empty <- true" would be blocked here.
			// and break statement would not be executed at once.

			// it needs to break this for-range at once.
			break
		}
	}
	os.Exit(0)

AppendError:
	errmsg := "epilog.End Error: Appender Append Failed"
	if err := lgr.appenders["stdio"].Append(errmsg); err != nil {
		fmt.Println("epilog.End Error: Appender stdio Append Failed")
	}
	os.Exit(0)
}

// Monitor is a method to monitor items in buffer
func (lgr *Logger) Monitor() {
	for {
		select {
		// if buffer is full, then process bufferItem in buffer.
		// flush the whole buffer until len(buffer.items) == 0.
		case _ = <-lgr.buffer.full:
			lgr.buffer.mtx.Lock()
			// use for-range to receive bufferItem.
			for item := range lgr.buffer.items {
				// Append of each appender should be called.
				// if err occur, ignore it and move on.
				for _, appender := range lgr.appenders {
					err := appender.Append(item.Serialize())
					if err != nil {
						errmsg := "epilog.Monitor Error: Appender Append Failed:"
						lgr.appenders["stdio"].Append(errmsg + err.Error())
					}
				}

				// if len(buffer.items) equal 0, then break for-range.
				if len(lgr.buffer.items) == 0 {
					lgr.buffer.empty <- true
					break
				}
			}
			lgr.buffer.mtx.Unlock()
		default:
		}
	}
}

// ChangeLevel is to change level
func (lgr *Logger) ChangeLevel(level logLevel) {
	lgr.mtx.Lock()
	lgr.level = level
	lgr.mtx.Unlock()
}

// RegisterAppender is to register appender
func (lgr *Logger) RegisterAppender(name string, appender Appender) (err error) {
	if _, ok := lgr.appenders[name]; ok {
		return errors.New("epilog.RegisterAppender error: appender already exists")
	}
	lgr.mtx.Lock()
	lgr.appenders[name] = appender
	lgr.mtx.Unlock()
	return nil
}

// Print is to print
func (lgr *Logger) Print(args ...interface{}) {
	innerPrint(lgr, false, NORMAL, args...)
	lgr.seqNum++
}

// Debugln is to debugln
func (lgr *Logger) Debugln(args ...interface{}) {
	if lgr.level <= DEBUG {
		innerPrint(lgr, true, DEBUG, args...)
		lgr.seqNum++
	}
}

// Infoln is to infoln
func (lgr *Logger) Infoln(args ...interface{}) {
	if lgr.level <= INFO {
		innerPrint(lgr, true, INFO, args...)
		lgr.seqNum++
	}
}

// Warningln is to warningln
func (lgr *Logger) Warningln(args ...interface{}) {
	if lgr.level <= WARNING {
		innerPrint(lgr, true, WARNING, args...)
		lgr.seqNum++
	}
}

// Fatalln is to fatalln
func (lgr *Logger) Fatalln(args ...interface{}) {
	innerPrint(lgr, true, FATAL, args...)
	os.Exit(0)
}
