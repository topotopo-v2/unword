package word

import (
	"time"

	"github.com/google/uuid"
)

type Word struct {
	ID            uuid.UUID `json:"id"`
	Word          string    `json:"word"`
	NativeScript  *string   `json:"native_script"`
	Pronunciation string    `json:"pronunciation"`
	Language      string    `json:"language"`
	Country       string    `json:"country"`
	CountryCode   string    `json:"country_code"`
	Definition    string    `json:"definition"`
	WordDate      time.Time `json:"word_date"`
	Source        *string   `json:"source"`
	CreatedAt     time.Time `json:"created_at"`
}
