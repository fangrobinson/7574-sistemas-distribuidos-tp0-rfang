package serialization

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
)

const SINGLE_BET = 1
const SINGLE_BET_ACK = 2
const MULTIPLE_BET = 3
const MULTIPLE_BET_ACK = 4
const NO_MORE_BETS = 5
const NO_MORE_BETS_ACK = 6
const GET_WINNERS = 7
const WAIT = 8
const WINNERS = 9

func IsVariableLength(code byte) bool {
	switch code {
	case byte(WINNERS):
		return true
	default:
		return false
	}
}

// GetMessageLength returns the length of the data associated with a specific message code.
func GetMessageLength(code byte) (int, error) {
	switch code {
	case byte(SINGLE_BET_ACK):
		return 0, nil
	case byte(MULTIPLE_BET_ACK):
		return 0, nil
	case byte(NO_MORE_BETS_ACK):
		return 0, nil
	case byte(WAIT):
		return 0, nil
	case byte(WINNERS):
		return 2, nil
	default:
		return 0, fmt.Errorf("unexpected message code: %v", code)
	}
}

// Extends a buffer with formatted bet
func BetExtendBuffer(b model.Bet, buffer *bytes.Buffer) error {
	buffer.WriteString(fmt.Sprintf("%-30s", b.Name)[:30])
	buffer.WriteString(fmt.Sprintf("%-20s", b.Surname)[:20])
	if err := binary.Write(buffer, binary.BigEndian, uint32(b.ID)); err != nil {
		return err
	}
	buffer.WriteString(fmt.Sprintf("%-10s", b.BirthDate)[:10])
	if err := binary.Write(buffer, binary.BigEndian, uint16(b.Number)); err != nil {
		return err
	}
	return nil
}

func EncodeBet(
	agencyId string,
	bet model.Bet,
) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	buffer.WriteByte(byte(SINGLE_BET))
	buffer.WriteString(fmt.Sprintf("%-3s", agencyId)[:3])
	err := BetExtendBuffer(bet, &buffer)
	if err != nil {
		return nil, err
	}
	return &buffer, nil
}

func EncodeMultipleBets(
	agencyId string,
	bets []model.Bet,
) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, uint8(MULTIPLE_BET)); err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("%-3s", agencyId)[:3])
	if err := binary.Write(&buffer, binary.BigEndian, uint8(len(bets))); err != nil {
		return nil, err
	}
	for _, bet := range bets {
		err := BetExtendBuffer(bet, &buffer)
		if err != nil {
			return nil, err
		}
	}
	return &buffer, nil
}

func EncodeNoMoreBets(
	agencyId string,
) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, uint8(NO_MORE_BETS)); err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("%-3s", agencyId)[:3])
	return &buffer, nil
}

func EncodeGetWinners(
	agencyId string,
) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, uint8(GET_WINNERS)); err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("%-3s", agencyId)[:3])
	return &buffer, nil
}
