package serialization

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
)

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
	buffer.WriteByte(byte(1))
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
	if err := binary.Write(&buffer, binary.BigEndian, uint8(3)); err != nil {
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
