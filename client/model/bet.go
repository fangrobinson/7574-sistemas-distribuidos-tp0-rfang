package model

import (
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
