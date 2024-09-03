package protocol

import (
	"net"
)

func ReceiveMessage(conn net.Conn) ([]byte, error) {
	code := make([]byte, 1)
	_, err := conn.Read(code)
	if err != nil {
		return nil, err
	}

	length, err := GetMessageLength(code[0])
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	if length > 0 {
		_, err := conn.Read(data)
		if err != nil {
			return nil, err
		}
	}
	message := append(code, data...)
	return message, nil
}
