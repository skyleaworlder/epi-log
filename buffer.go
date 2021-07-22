package epilog

import (
	"fmt"
	"sync"
	"time"
)

// Buffer is a struct
type Buffer struct {
	maxItemNum  int
	items       chan Item
	empty, full chan bool
	mtx         *sync.Mutex
}

// BufferItem is a struct
type BufferItem struct {
	id       int
	typ      string
	time     time.Time
	filename string
	line     int
	content  string
}

// NewBuffer is to create a Buffer
func NewBuffer(maxItemNum int) (b *Buffer) {
	return &Buffer{
		maxItemNum: maxItemNum,
		items:      make(chan Item, maxItemNum),
		empty:      make(chan bool),
		full:       make(chan bool),
		mtx:        &sync.Mutex{},
	}
}

// Put is to push bufferitem into buffer.
// if buffer.items is full, then buffer.full <- true.
// buffer.full <- true will let Logger.Monitor process full buffer.
func (b *Buffer) Put(item Item) {
	// if b.items is full, b.full <- true
	select {
	case b.items <- item:
	default:
		b.full <- true

	SpinLock:
		for {
			select {
			// Monitor: mgr.buffer.empty <- true
			case <-b.empty:
				// if b.items <- item deleted, item will be missing.
				// use break label to break a for-select statement.
				b.items <- item
				break SpinLock
			default:
			}
		}
	}
}

// Serialize is to implement interface "Item"
func (i BufferItem) Serialize() (res string) {
	res = fmt.Sprintf(
		"[%d] %s %s: %s Line_%d -> %s",
		i.id, i.time, i.typ, i.filename, i.line, i.content,
	)
	return
}
