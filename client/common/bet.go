package common

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/spf13/viper"
)

type Bet struct {
	Name      string
	Surname   string
	ID        int
	BirthDate string
	Number    int
}

// Initializes a new Bet receiving the configuration
// as a parameter
func NewBet(v *viper.Viper) *Bet {
	b := &Bet{
		Name:      v.GetString("nombre"),
		Surname:   v.GetString("apellido"),
		ID:        v.GetInt("documento"),
		BirthDate: v.GetString("nacimiento"),
		Number:    v.GetInt("numero"),
	}
	return b
}

// Extends a buffer with formatted bet
func (b *Bet) extendBuffer(buffer *bytes.Buffer) error {
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
