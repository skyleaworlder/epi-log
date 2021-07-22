package epilog

import (
	"errors"
	"fmt"
	"net"
	"os"
)

// StdAppender is an appender.
type StdAppender struct {
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
