package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/topotopo-v2/unword/internal/word"
)

type seedWord struct {
	Word          string  `json:"word"`
	NativeScript  *string `json:"native_script"`
	Pronunciation string  `json:"pronunciation"`
	Language      string  `json:"language"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Definition    string  `json:"definition"`
	WordDate      string  `json:"word_date"`
	Source        *string `json:"source"`
}

func main() {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	data, err := os.ReadFile("data/words.json")
	if err != nil {
		log.Fatal(err)
	}

	var seedWords []seedWord

	if err := json.Unmarshal(data, &seedWords); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repository := word.NewRepository(pool)

	for _, item := range seedWords {
		wordDate, err := time.Parse("2006-01-02", item.WordDate)
		if err != nil {
			log.Fatalf("invalid date %q: %v", item.WordDate, err)
		}

		w := word.Word{
			ID:            uuid.New(),
			Word:          item.Word,
			NativeScript:  item.NativeScript,
			Pronunciation: item.Pronunciation,
			Language:      item.Language,
			Country:       item.Country,
			CountryCode:   item.CountryCode,
			Definition:    item.Definition,
			WordDate:      wordDate,
			Source:        item.Source,
		}

		if err := repository.Create(ctx, w); err != nil {
			log.Fatalf(
				"failed to insert %q for %s: %v",
				item.Word,
				item.WordDate,
				err,
			)
		}
	}

	fmt.Printf("Inserted %d words\n", len(seedWords))
}