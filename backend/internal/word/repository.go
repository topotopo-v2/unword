package word

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Repository struct {
	pool DB
}

func NewRepository(pool DB) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Create(ctx context.Context, w Word) error {
	_, err := r.pool.Exec(
		ctx,
		`
		INSERT INTO words (
			id,
			word,
			native_script,
			pronunciation,
			language,
			country,
			country_code,
			definition,
			word_date,
			source
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)
		`,
		w.ID,
		w.Word,
		w.NativeScript,
		w.Pronunciation,
		w.Language,
		w.Country,
		w.CountryCode,
		w.Definition,
		w.WordDate,
		w.Source,
	)

	return err
}

func (r *Repository) GetByDate(
	ctx context.Context,
	date time.Time,
) (*Word, error) {
	var w Word

	err := r.pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			word,
			native_script,
			pronunciation,
			language,
			country,
			country_code,
			definition,
			word_date,
			source,
			created_at
		FROM words
		WHERE word_date = $1
		`,
		date,
	).Scan(
		&w.ID,
		&w.Word,
		&w.NativeScript,
		&w.Pronunciation,
		&w.Language,
		&w.Country,
		&w.CountryCode,
		&w.Definition,
		&w.WordDate,
		&w.Source,
		&w.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *Repository) GetWordsByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]Word, error) {
	rows, err := r.pool.Query(
		ctx,
		`
        SELECT
            id,
            word,
            native_script,
            pronunciation,
            language,
            country,
            country_code,
            definition,
            word_date,
            source,
            created_at
        FROM words
        WHERE id = ANY($1)
        ORDER BY word_date DESC
        `,
		ids,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	words := make([]Word, 0)

	for rows.Next() {
		var w Word

		err := rows.Scan(
			&w.ID,
			&w.Word,
			&w.NativeScript,
			&w.Pronunciation,
			&w.Language,
			&w.Country,
			&w.CountryCode,
			&w.Definition,
			&w.WordDate,
			&w.Source,
			&w.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		words = append(words, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
