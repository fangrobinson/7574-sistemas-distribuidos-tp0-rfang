package protocol

import (
	"errors"
	"net"
)

func SendMessage(conn net.Conn, data []byte) error {
	sent := 0
	for sent < len(data) {
		n, err := conn.Write(data[sent:])
		if err != nil {
			return err
		}
		if n == 0 {
			return errors.New("connection was lost")
		}
		sent += n
	}
	return nil
}
