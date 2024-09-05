package protocol

import (
	"encoding/binary"
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/serialization"
)

func RecvFixed(conn net.Conn, code byte) ([]byte, error) {
	length, err := serialization.GetMessageLength(code)
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
	return data, nil
}

func RecvVariable(conn net.Conn, code byte) ([]byte, error) {
	amountLength := 2
	winnersAmount := make([]byte, amountLength)
	totalRead := 0
	for totalRead < amountLength {
		n, err := conn.Read(winnersAmount[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}
	winnersAmountInt := int(binary.BigEndian.Uint16(winnersAmount))
	winnersLength := winnersAmountInt * 4
	winners := make([]byte, winnersLength)
	totalRead = 0
	for totalRead < winnersLength {
		n, err := conn.Read(winners[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}
	message := append(winnersAmount, winners...)
	return message, nil
}

func ReceiveMessage(conn net.Conn) ([]byte, error) {
	code := make([]byte, 1)
	_, err := conn.Read(code)
	if err != nil {
		return nil, err
	}

	if serialization.IsVariableLength(code[0]) {
		data, err := RecvVariable(conn, code[0])
		if err != nil {
			return nil, err
		}
		message := append(code, data...)
		return message, nil
	} else {
		data, err := RecvFixed(conn, code[0])
		if err != nil {
			return nil, err
		}
		message := append(code, data...)
		return message, nil
	}
}

func ReceiveWinners(conn net.Conn) ([]int, error) {
	code := make([]byte, 1)
	_, err := conn.Read(code)
	if err != nil {
		return nil, err
	}
	amountLength := 2
	winnersAmount := make([]byte, amountLength)
	totalRead := 0
	for totalRead < amountLength {
		n, err := conn.Read(winnersAmount[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}
	winnersAmountInt := int(binary.BigEndian.Uint16(winnersAmount))
	winnersLength := winnersAmountInt * 4
	winners := make([]byte, winnersLength)
	totalRead = 0
	for totalRead < winnersLength {
		n, err := conn.Read(winners[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}

	winnersIDs := make([]int, winnersAmountInt)
	for i := 0; i < winnersAmountInt; i++ {
		idStart := i * 4
		idEnd := idStart + 4
		winnerId := int(binary.BigEndian.Uint32(winners[idStart:idEnd]))
		winnersIDs = append(winnersIDs, winnerId)
	}

	return winnersIDs, nil
}
