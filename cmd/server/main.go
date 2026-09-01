package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/topotopo-v2/unword/internal/word"
	"os"
)

func main() {
	// Get environment
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Establish connection to PostgresSQL
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("successfully connected to PostgreSQL")

	wordRepository := word.NewRepository(pool)
	wordHandler := word.NewHandler(wordRepository)

	http.HandleFunc(
		"/api/words/today",
		wordHandler.GetToday,
	)

/* 	// Test insert word
		wordID := uuid.New()
		testWord := word.Word{
			ID:            wordID,
			Word:          "ubuntu",
			Pronunciation: "oo-BOON-too",
			Language:      "Nguni",
			Country:       "South Africa",
			CountryCode:   "ZA",
			Definition:    "A concept associated with humanity, interconnectedness, and community.",
			WordDate:      time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC),
		}

		err = wordRepository.Create(ctx, testWord)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("word inserted successfully") */

/* 	// Test get word
		result, err := wordRepository.GetByDate(
			ctx,
			time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC),
		)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("found word: %s", result.Word) */

	http.Handle("/health", healthHandler(pool))
	log.Println("server listening on :8080")

	server := http.ListenAndServe(":8080", nil)
	if server != nil {
		log.Fatal(server)
	}
}

func healthHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := pool.Ping(r.Context()); err != nil {
			// PostgreSQL is unavailable
			http.Error(w, `{"status":"unavailable"}`, http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))

	}
}
