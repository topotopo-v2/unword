package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// Get environment
	err := godotenv.Load(".env") // connected just fine because im running inside the root project?
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("successfully connected to PostgreSQL")

	http.Handle("/health", healthHandler(pool))
	log.Println("server listening on :8080")

	serverr := http.ListenAndServe(":8080", nil)
	if serverr != nil {
		log.Fatal(serverr)
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
