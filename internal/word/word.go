package word

import (
	"time"

	"github.com/google/uuid"
)

type Word struct {
	ID            uuid.UUID
	Word          string
	NativeScript  *string
	Pronunciation string
	Language      string
	Country       string
	CountryCode   string
	Definition    string
	WordDate      time.Time
	Source        *string
	CreatedAt     time.Time
}
