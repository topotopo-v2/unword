package word

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const insertWordSQL = `INSERT INTO words
	( id, word, native_script, pronunciation, language, country, country_code, definition, word_date, source )
	VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 )`

const selectWordColumns = `SELECT id, word, native_script, pronunciation, language,
	country, country_code, definition, word_date, source, created_at FROM words`

var wordColumns = []string{
	"id", "word", "native_script", "pronunciation", "language", "country",
	"country_code", "definition", "word_date", "source", "created_at",
}

func repositoryTestWord() Word {
	nativeScript, source := "saudade", "https://example.com/saudade"
	return Word{
		ID:   uuid.MustParse("11111111-1111-4111-8111-111111111111"),
		Word: "saudade", NativeScript: &nativeScript, Pronunciation: "sow-dah-deh",
		Language: "Portuguese", Country: "Portugal", CountryCode: "PT",
		Definition: "A feeling of longing.", Source: &source,
		WordDate:  time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Date(2026, 9, 10, 12, 30, 0, 0, time.UTC),
	}
}

func wordRow(w Word) []any {
	return []any{w.ID, w.Word, w.NativeScript, w.Pronunciation, w.Language,
		w.Country, w.CountryCode, w.Definition, w.WordDate, w.Source, w.CreatedAt}
}

func newRepositoryMock(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, mock.ExpectationsWereMet(), "database expectations")
		mock.Close()
	})
	return mock
}

func TestRepositoryCreate(t *testing.T) {
	word := repositoryTestWord()
	withoutOptional := word
	withoutOptional.NativeScript, withoutOptional.Source = nil, nil
	tests := []struct {
		name  string
		word  Word
		dbErr error
	}{
		{name: "success with optional fields", word: word},
		{name: "success without optional fields", word: withoutOptional},
		{name: "database failure", word: word, dbErr: errors.New("database unavailable")},
		{name: "duplicate date", word: word, dbErr: &pgconn.PgError{Code: "23505", ConstraintName: "words_word_date_key"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newRepositoryMock(t)
			w := tt.word
			expect := mock.ExpectExec(insertWordSQL).WithArgs(w.ID, w.Word, w.NativeScript,
				w.Pronunciation, w.Language, w.Country, w.CountryCode, w.Definition, w.WordDate, w.Source)
			if tt.dbErr != nil {
				expect.WillReturnError(tt.dbErr)
			} else {
				expect.WillReturnResult(pgxmock.NewResult("INSERT", 1))
			}
			err := NewRepository(mock).Create(context.Background(), w)
			if tt.dbErr != nil {
				require.ErrorIs(t, err, tt.dbErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRepositoryGetByDate(t *testing.T) {
	word := repositoryTestWord()
	withoutOptional := word
	withoutOptional.NativeScript, withoutOptional.Source = nil, nil
	tests := []struct {
		name         string
		word         Word
		noRows       bool
		nullOptional bool
		scanFailure  bool
		dbErr        error
	}{
		{name: "word found", word: word},
		{name: "NULL optional fields", word: withoutOptional, nullOptional: true},
		{name: "no matching row", noRows: true},
		{name: "database failure", dbErr: errors.New("query failed")},
		{name: "scan failure", word: word, scanFailure: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newRepositoryMock(t)
			expect := mock.ExpectQuery(selectWordColumns + " WHERE word_date = $1").WithArgs(word.WordDate)
			if tt.dbErr != nil {
				expect.WillReturnError(tt.dbErr)
			} else {
				rows := pgxmock.NewRows(wordColumns)
				if !tt.noRows {
					values := wordRow(tt.word)
					if tt.nullOptional {
						values[2], values[9] = nil, nil
					}
					if tt.scanFailure {
						values[8] = struct{}{}
					}
					rows.AddRow(values...)
				}
				expect.WillReturnRows(rows)
			}
			got, err := NewRepository(mock).GetByDate(context.Background(), word.WordDate)
			switch {
			case tt.dbErr != nil:
				require.ErrorIs(t, err, tt.dbErr)
				assert.Nil(t, got)
			case tt.noRows:
				require.ErrorIs(t, err, pgx.ErrNoRows)
				assert.Nil(t, got)
			case tt.scanFailure:
				require.Error(t, err)
				assert.Contains(t, err.Error(), "word_date")
				assert.Nil(t, got)
			default:
				require.NoError(t, err)
				assert.Equal(t, &tt.word, got)
			}
		})
	}
}

func TestRepositoryGetWordsByIDs(t *testing.T) {
	first := repositoryTestWord()
	second := first
	second.ID = uuid.MustParse("22222222-2222-4222-8222-222222222222")
	second.Word, second.Language, second.Country, second.CountryCode = "ubuntu", "Nguni", "South Africa", "ZA"
	second.WordDate = first.WordDate.AddDate(0, 0, -1)
	second.NativeScript, second.Source = nil, nil
	ids := []uuid.UUID{second.ID, first.ID}
	tests := []struct {
		name        string
		ids         []uuid.UUID
		words       []Word
		dbErr       error
		scanFailure bool
		endErr      error
	}{
		{name: "multiple rows including NULL optional fields", ids: ids, words: []Word{first, second}},
		{name: "no matching rows", ids: ids, words: []Word{}},
		{name: "empty IDs", ids: []uuid.UUID{}, words: []Word{}},
		{name: "query failure", ids: ids, dbErr: errors.New("query failed")},
		{name: "scan failure discards partial results", ids: ids, words: []Word{first, second}, scanFailure: true},
		{name: "end of rows error discards partial results", ids: ids, words: []Word{first}, endErr: errors.New("row stream failed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newRepositoryMock(t)
			expect := mock.ExpectQuery(selectWordColumns + " WHERE id = ANY($1) ORDER BY word_date DESC").WithArgs(tt.ids)
			if tt.dbErr != nil {
				expect.WillReturnError(tt.dbErr)
			} else {
				rows := pgxmock.NewRows(wordColumns)
				for i, w := range tt.words {
					values := wordRow(w)
					if w.NativeScript == nil {
						values[2] = nil
					}
					if w.Source == nil {
						values[9] = nil
					}
					if tt.scanFailure && i == 1 {
						values[8] = struct{}{}
					}
					rows.AddRow(values...)
				}
				if tt.endErr != nil {
					// pgxmock exposes CloseError through Err after Next reaches the end.
					rows.CloseError(tt.endErr)
				}
				expect.WillReturnRows(rows).RowsWillBeClosed()
			}
			got, err := NewRepository(mock).GetWordsByIDs(context.Background(), tt.ids)
			switch {
			case tt.dbErr != nil:
				require.ErrorIs(t, err, tt.dbErr)
				assert.Nil(t, got)
			case tt.scanFailure:
				require.Error(t, err)
				assert.Contains(t, err.Error(), "word_date")
				assert.Nil(t, got)
			case tt.endErr != nil:
				require.ErrorIs(t, err, tt.endErr)
				assert.Nil(t, got)
			default:
				require.NoError(t, err)
				assert.Equal(t, tt.words, got)
			}
		})
	}
}
