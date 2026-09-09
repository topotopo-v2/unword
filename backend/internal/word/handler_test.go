package word

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
)

func TestValidateWord(t *testing.T) {
	tests := []struct {
		name    string
		word    Word
		wantErr bool
	}{
		{
			name: "valid word",
			word: Word{
				ID:            uuid.New(),
				Word:          "saudade",
				Pronunciation: "sow-dah-deh",
				Language:      "Portuguese",
				Country:       "Portugal",
				CountryCode:   "PT",
				Definition:    "A deep emotional state...",
				WordDate:      time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing word",
			word: Word{
				Language:   "Portuguese",
				Definition: "A definition",
				WordDate:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing language",
			word: Word{
				Word:       "saudade",
				Definition: "A definition",
				WordDate:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing definition",
			word: Word{
				Word:     "saudade",
				Language: "Portuguese",
				WordDate: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid country code",
			word: Word{
				Word:        "saudade",
				Language:    "Portuguese",
				CountryCode: "POR",
				Definition:  "A definition",
				WordDate:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing word date",
			word: Word{
				Word:        "saudade",
				Language:    "Portuguese",
				CountryCode: "PT",
				Definition:  "A definition",
			},
			wantErr: true,
		},
		{
			name: "blank word",
			word: Word{
				Word:        "   ",
				Language:    "Portuguese",
				CountryCode: "PT",
				Definition:  "A definition",
				WordDate:    time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWord(tt.word)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateWord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateWord_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/words",
		strings.NewReader(`{"word":`),
	)

	rec := httptest.NewRecorder()

	handler := &Handler{}

	handler.CreateWord(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
