package epilog

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// Logger is a struct
type Logger struct {
	level logLevel

	appenders map[string]Appender
	buffer    *Buffer

	cnt int

	mtx *sync.Mutex
}

// New is a constructor
func New(level logLevel, maxItemNum int) (lgr *Logger) {
	return &Logger{
		level:  level,
		buffer: NewBuffer(maxItemNum),
		cnt:    0,
		mtx:    &sync.Mutex{},
	}
}

// Use is to use plugins
func (lgr *Logger) Use() {
}

// Run is to run
func (lgr *Logger) Run() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		for sig := range sigs {
			switch sig {
			case syscall.SIGINT:
				fmt.Println("syscall.SIGINT:", sig)
				lgr.UrgentExit()
			case syscall.SIGTERM:
				fmt.Println("syscall.SIGTERM:", sig)
				lgr.UrgentExit()
			case syscall.SIGKILL:
				fmt.Println("syscall.SIGKILL:", sig)
				lgr.UrgentExit()
			default:
				fmt.Println("what?")
			}
		}
	}()

	lgr.Monitor()
}

// UrgentExit is to exit
func (lgr *Logger) UrgentExit() {
	os.Exit(0)
}

// Monitor is a method to monitor items in buffer
func (lgr *Logger) Monitor() {
	for {
		select {
		case _ = <-lgr.buffer.full:
			lgr.buffer.mtx.Lock()
			for item := range lgr.buffer.items {
				// TODO: send/write items here!!!
				fmt.Println(item.Serialize())
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
	// Print ignore Logger.level
	fmt.Print(args...)
	content := fmt.Sprint(args...)
	_, filename, line, _ := runtime.Caller(0)

	// put into buffer
	lgr.buffer.Put(BufferItem{
		id: lgr.cnt, time: time.Now(),
		filename: filename, line: line,
		content: content,
	})

	lgr.cnt++
}
