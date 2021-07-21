package epilog

import (
	"errors"
	"fmt"
	"sync"
)

// Logger is a struct
type Logger struct {
	level     logLevel
	appenders map[string]Appender
	buffer    *Buffer
	mtx       *sync.Mutex
}

// New is a constructor
func New(level logLevel, maxItemNum int) (mgr *Logger) {
	return &Logger{
		level:  level,
		buffer: NewBuffer(maxItemNum),
		mtx:    &sync.Mutex{},
	}
}

// Use is to use plugins
func (mgr *Logger) Use() {
}

// Monitor is a method to monitor items in buffer
func (mgr *Logger) Monitor() {
	for {
		select {
		case _ = <-mgr.buffer.full:
			mgr.buffer.mtx.Lock()
			for item := range mgr.buffer.items {
				// TODO: send/write items here!!!
				fmt.Println(item.Serialize())
				if len(mgr.buffer.items) == 0 {
					mgr.buffer.empty <- true
					break
				}
			}
			mgr.buffer.mtx.Unlock()
		default:
		}
	}
}

// ChangeLevel is to change level
func (mgr *Logger) ChangeLevel(level logLevel) {
	mgr.mtx.Lock()
	mgr.level = level
	mgr.mtx.Unlock()
}

// RegisterAppender is to register appender
func (mgr *Logger) RegisterAppender(name string, appender Appender) (err error) {
	if _, ok := mgr.appenders[name]; ok {
		return errors.New("epilog.RegisterAppender error: appender already exists")
	}
	mgr.mtx.Lock()
	mgr.appenders[name] = appender
	mgr.mtx.Unlock()
	return nil
}
