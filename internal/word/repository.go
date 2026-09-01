package word

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
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
