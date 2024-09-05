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
		totalRead := 0
		for totalRead < length {
			n, err := conn.Read(data[totalRead:])
			if err != nil {
				return nil, err
			}
			totalRead += n
		}
	}
	message := append(code, data...)
	return message, nil
}
