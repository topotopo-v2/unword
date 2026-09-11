package word

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
)

type fakeWordRepository struct {
	createFn        func(context.Context, Word) error
	getByDateFn     func(context.Context, time.Time) (*Word, error)
	getWordsByIDsFn func(context.Context, []uuid.UUID) ([]Word, error)
}

func (f *fakeWordRepository) Create(ctx context.Context, word Word) error {
	return f.createFn(ctx, word)
}

func (f *fakeWordRepository) GetByDate(ctx context.Context, date time.Time) (*Word, error) {
	return f.getByDateFn(ctx, date)
}

func (f *fakeWordRepository) GetWordsByIDs(ctx context.Context, ids []uuid.UUID) ([]Word, error) {
	return f.getWordsByIDsFn(ctx, ids)
}

func TestGetToday(t *testing.T) {
	word := &Word{ID: uuid.New(), Word: "saudade", CountryCode: "PT"}
	tests := []struct {
		name       string
		timezone   string
		repoErr    error
		wantStatus int
		wantError  string
		wantCalls  int
	}{
		{name: "missing timezone", wantStatus: http.StatusBadRequest, wantError: "timezone is required"},
		{name: "invalid timezone", timezone: "Invalid/Zone", wantStatus: http.StatusBadRequest, wantError: "invalid timezone"},
		{name: "word not found", timezone: "UTC", repoErr: pgx.ErrNoRows, wantStatus: http.StatusNotFound, wantError: "word not found", wantCalls: 1},
		{name: "wrapped not found", timezone: "UTC", repoErr: fmt.Errorf("lookup: %w", pgx.ErrNoRows), wantStatus: http.StatusNotFound, wantError: "word not found", wantCalls: 1},
		{name: "repository failure", timezone: "UTC", repoErr: errors.New("internal database detail"), wantStatus: http.StatusInternalServerError, wantError: "failed to get today's word", wantCalls: 1},
		{name: "success UTC", timezone: "UTC", wantStatus: http.StatusOK, wantCalls: 1},
		{name: "success east of UTC", timezone: "Pacific/Kiritimati", wantStatus: http.StatusOK, wantCalls: 1},
		{name: "success west of UTC", timezone: "Pacific/Pago_Pago", wantStatus: http.StatusOK, wantCalls: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			var receivedDate time.Time
			var receivedContext context.Context
			repo := &fakeWordRepository{getByDateFn: func(ctx context.Context, date time.Time) (*Word, error) {
				calls++
				receivedDate, receivedContext = date, ctx
				return word, tt.repoErr
			}}
			req := httptest.NewRequest(http.MethodGet, "/api/words/today?timezone="+url.QueryEscape(tt.timezone), nil)
			rec := httptest.NewRecorder()
			before := time.Now()
			NewHandler(repo).GetToday(rec, req)
			after := time.Now()
			require.Equal(t, tt.wantStatus, rec.Code, "response body: %s", rec.Body.String())
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			require.Equal(t, tt.wantCalls, calls)
			if calls == 1 {
				assert.Equal(t, req.Context(), receivedContext)
				location, err := time.LoadLocation(tt.timezone)
				require.NoError(t, err)
				// Allow either date if the request crosses local midnight.
				assert.Contains(t, []string{before.In(location).Format(time.DateOnly), after.In(location).Format(time.DateOnly)}, receivedDate.Format(time.DateOnly))
				assert.Equal(t, time.UTC, receivedDate.Location())
				assert.Equal(t, "00:00:00.000000000", receivedDate.Format("15:04:05.000000000"))
			}
			if tt.wantError != "" {
				var response map[string]string
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
				assert.Equal(t, map[string]string{"error": tt.wantError}, response)
				return
			}
			var response Word
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			assert.Equal(t, *word, response)
		})
	}
}

func TestGetWords(t *testing.T) {
	firstID, secondID := uuid.New(), uuid.New()
	words := []Word{{ID: firstID, Word: "saudade"}, {ID: secondID, Word: "ubuntu"}}
	tests := []struct {
		name          string
		ids           string
		repoWords     []Word
		repoErr       error
		wantStatus    int
		wantError     string
		wantPlainText bool
		wantIDs       []uuid.UUID
	}{
		{name: "missing IDs", wantStatus: http.StatusBadRequest, wantError: "ids is required"},
		{name: "invalid ID", ids: "invalid", wantStatus: http.StatusBadRequest, wantError: "invalid word id", wantPlainText: true},
		{name: "invalid ID after valid ID", ids: firstID.String() + ",invalid", wantStatus: http.StatusBadRequest, wantError: "invalid word id", wantPlainText: true},
		{name: "empty ID in list", ids: firstID.String() + ",", wantStatus: http.StatusBadRequest, wantError: "invalid word id", wantPlainText: true},
		{name: "repository failure", ids: firstID.String(), repoErr: errors.New("internal database detail"), wantStatus: http.StatusInternalServerError, wantError: "failed to get words", wantPlainText: true, wantIDs: []uuid.UUID{firstID}},
		{name: "single word", ids: firstID.String(), repoWords: words[:1], wantStatus: http.StatusOK, wantIDs: []uuid.UUID{firstID}},
		{name: "multiple words", ids: firstID.String() + "," + secondID.String(), repoWords: words, wantStatus: http.StatusOK, wantIDs: []uuid.UUID{firstID, secondID}},
		{name: "no matching words", ids: firstID.String(), repoWords: []Word{}, wantStatus: http.StatusOK, wantIDs: []uuid.UUID{firstID}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			var receivedIDs []uuid.UUID
			var receivedContext context.Context
			repo := &fakeWordRepository{getWordsByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]Word, error) {
				calls++
				receivedIDs, receivedContext = ids, ctx
				return tt.repoWords, tt.repoErr
			}}
			req := httptest.NewRequest(http.MethodGet, "/api/words?ids="+url.QueryEscape(tt.ids), nil)
			rec := httptest.NewRecorder()
			NewHandler(repo).GetWords(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code, "response body: %s", rec.Body.String())
			if tt.wantIDs == nil {
				require.Zero(t, calls)
			} else {
				require.Equal(t, 1, calls)
				assert.Equal(t, tt.wantIDs, receivedIDs)
				assert.Equal(t, req.Context(), receivedContext)
			}
			if tt.wantPlainText {
				assert.Equal(t, "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
				assert.Equal(t, tt.wantError+"\n", rec.Body.String())
				return
			}
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			if tt.wantError != "" {
				var response map[string]string
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
				assert.Equal(t, map[string]string{"error": tt.wantError}, response)
				return
			}
			var response []Word
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			assert.Equal(t, tt.repoWords, response)
		})
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		value    any
		wantJSON string
	}{
		{name: "object", status: http.StatusCreated, value: map[string]string{"word": "saudade"}, wantJSON: `{"word":"saudade"}`},
		{name: "list", status: http.StatusOK, value: []string{"saudade", "ubuntu"}, wantJSON: `["saudade","ubuntu"]`},
		{name: "empty list", status: http.StatusOK, value: []string{}, wantJSON: `[]`},
		{name: "nil", status: http.StatusOK, value: nil, wantJSON: `null`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeJSON(rec, tt.status, tt.value)
			require.Equal(t, tt.status, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantJSON, rec.Body.String())
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		message string
	}{
		{name: "bad request", status: http.StatusBadRequest, message: "invalid request body"},
		{name: "server error", status: http.StatusInternalServerError, message: "failed to create word"},
		{name: "special characters", status: http.StatusBadRequest, message: "invalid \"word\"\ntry again: café"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeError(rec, tt.message, tt.status)
			require.Equal(t, tt.status, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			var response map[string]string
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			assert.Equal(t, map[string]string{"error": tt.message}, response)
		})
	}
}

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

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateWord(t *testing.T) {
	validWord := Word{
		Word:          "saudade",
		Pronunciation: "sow-dah-deh",
		Language:      "Portuguese",
		Country:       "Portugal",
		CountryCode:   " pt ",
		Definition:    "A feeling of longing.",
		WordDate:      time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC),
	}
	invalidWord := validWord
	invalidWord.Word = "   "
	wordWithID := validWord
	wordWithID.CountryCode = "PT"
	wordWithID.ID = uuid.MustParse("11111111-1111-4111-8111-111111111111")

	tests := []struct {
		name            string
		word            Word
		rawBody         string
		repoErr         error
		wantStatus      int
		wantError       string
		wantCalls       int
		wantCountryCode string
	}{
		{
			name: "invalid JSON", rawBody: `{"word":`,
			wantStatus: http.StatusBadRequest, wantError: "invalid request body",
		},
		{
			name: "invalid word", word: invalidWord,
			wantStatus: http.StatusBadRequest, wantError: "word cannot be blank",
		},
		{
			name: "duplicate date", word: validWord,
			repoErr:    &pgconn.PgError{Code: "23505", Message: "internal database detail"},
			wantStatus: http.StatusConflict, wantError: "word already exists for this date", wantCalls: 1,
		},
		{
			name: "repository failure", word: validWord,
			repoErr:    errors.New("internal database detail"),
			wantStatus: http.StatusInternalServerError, wantError: "failed to create word", wantCalls: 1,
		},
		{
			name: "success generates ID and normalizes country code", word: validWord,
			wantCountryCode: "PT",
			wantStatus:      http.StatusCreated, wantCalls: 1,
		},
		{
			name: "success preserves provided ID", word: wordWithID,
			wantCountryCode: "PT",
			wantStatus:      http.StatusCreated, wantCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := tt.rawBody
			if body == "" {
				data, err := json.Marshal(tt.word)
				require.NoError(t, err, "encode request body")
				body = string(data)
			}

			var received Word
			calls := 0
			repo := &fakeWordRepository{
				createFn: func(ctx context.Context, word Word) error {
					calls++
					received = word
					return tt.repoErr
				},
			}
			req := httptest.NewRequest(http.MethodPost, "/api/words", strings.NewReader(body))
			rec := httptest.NewRecorder()
			NewHandler(repo).CreateWord(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code, "response body: %s", rec.Body.String())
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			require.Equal(t, tt.wantCalls, calls, "repository call count")

			if tt.wantError != "" {
				var response map[string]string
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response), "decode error response")
				assert.Equal(t, map[string]string{"error": tt.wantError}, response)
				return
			}

			wantWord := tt.word
			wantWord.CountryCode = tt.wantCountryCode
			assert.Equal(t, tt.wantCountryCode, received.CountryCode, "country code")
			if wantWord.ID == uuid.Nil {
				require.NotEqual(t, uuid.Nil, received.ID, "expected a generated non-nil UUID")
				wantWord.ID = received.ID
			}
			assert.Equal(t, wantWord, received, "word passed to repository")

			var response Word
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response), "decode success response")
			assert.Equal(t, wantWord, response, "response word")
		})
	}
}
