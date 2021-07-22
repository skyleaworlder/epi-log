package epilog

import (
	"errors"
	"fmt"
	"net"
	"os"
)

var (
	appenderID int = 0
)

// StdAppender is an appender.
type StdAppender struct {
	id int
}

// NewStdAppender is default constructor of appender.
func NewStdAppender() *StdAppender {
	defer func() { appenderID++ }()
	return &StdAppender{id: appenderID}
}

// Append is to implement interface "Appender"
func (std *StdAppender) Append(content string) (err error) {
	_, err = fmt.Println(content)
	return
}

// FileAppender is an appender.
// name is the key of appenders map.
type FileAppender struct {
	id   int
	name string
	fd   *os.File
}

// NewFileAppender is a constructor of FileAppender.
func NewFileAppender(name string, fd *os.File) *FileAppender {
	defer func() { appenderID++ }()
	return &FileAppender{
		id:   appenderID,
		name: name,
		fd:   fd,
	}
}

// Append is to implement interface "Appender"
func (fa *FileAppender) Append(content string) (err error) {
	buf := []byte(content)
	if n, err := fa.fd.Write(buf); n < len(content) || err != nil {
		errmsg := "epilog.Append error: FileAppender Write Failed"
		return errors.New(errmsg + err.Error())
	}
	return
}

// SocketAppender is an appender.
// name is the key of appender map.
type SocketAppender struct {
	id   int
	name string
	ip   net.TCPAddr
}

// NewSocketAppender is a constructor of SocketAppender.
func NewSocketAppender(name string, ip net.TCPAddr) *SocketAppender {
	defer func() { appenderID++ }()
	return &SocketAppender{
		id:   appenderID,
		name: name,
		ip:   ip,
	}
}

// Append is to implement interface "Appender"
func (sa *SocketAppender) Append(content string) (err error) {
	conn, err := net.DialTCP("tcp4", nil, &sa.ip)
	if err != nil {
		errmsg := "epilog.Append error: SocketAppender Connect Failed"
		return errors.New(errmsg + err.Error())
	}

	buf := []byte(content)
	if n, err := conn.Write(buf); n < len(content) || err != nil {
		errmsg := "epilog.Append error: SocketAppender Write Failed"
		return errors.New(errmsg + err.Error())
	}
	return
}
