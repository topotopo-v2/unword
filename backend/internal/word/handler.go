package word

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) Words(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetWords(w, r)

	case http.MethodPost:
		h.CreateWord(w, r)

	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) CreateWord(w http.ResponseWriter, r *http.Request) {
	var word Word

	err := json.NewDecoder(r.Body).Decode(&word)
	if err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validateWord(word); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert country code to uppercase
	word.CountryCode = strings.ToUpper(strings.TrimSpace(word.CountryCode))

	if word.ID == uuid.Nil {
		word.ID = uuid.New()
	}

	err = h.repository.Create(r.Context(), word)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, "word already exists for this date", http.StatusConflict)
			return
		}

		writeError(w, "failed to create word", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(word)
}

func (h *Handler) GetToday(w http.ResponseWriter, r *http.Request) {
	timezone := r.URL.Query().Get("timezone")

	if timezone == "" {
		writeError(w, "timezone is required", http.StatusBadRequest)
		return
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		writeError(w, "invalid timezone", http.StatusBadRequest)
		return
	}

	now := time.Now().In(location)

	today := time.Date(
        now.Year(),
        now.Month(),
        now.Day(),
        0, 0, 0, 0,
        time.UTC,
    )

	word, err := h.repository.GetByDate(
		r.Context(),
		today,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, "word not found", http.StatusNotFound)
			return
		}
		writeError(w, "failed to get today's word", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, word)
}

func (h *Handler) GetWords(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")

	if idsParam == "" {
		writeError(w, "ids is required", http.StatusBadRequest)
		return
	}

	idStrings := strings.Split(idsParam, ",")

	ids := make([]uuid.UUID, 0, len(idStrings))

	for _, idString := range idStrings {
		id, err := uuid.Parse(idString)
		if err != nil {
			http.Error(w, "invalid word id", http.StatusBadRequest)
			return
		}

		ids = append(ids, id)
	}

	words, err := h.repository.GetWordsByIDs(r.Context(), ids)
	if err != nil {
		http.Error(w, "failed to get words", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, words)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func validateWord(word Word) error {
	if strings.TrimSpace(word.Word) == "" {
		return errors.New("word cannot be blank")
	}

	if strings.TrimSpace(word.Language) == "" {
		return errors.New("language cannot be blank")
	}

	if strings.TrimSpace(word.Definition) == "" {
		return errors.New("definition cannot be blank")
	}

	if len(strings.TrimSpace(word.CountryCode)) != 2 {
		return errors.New("country_code must be 2 characters")
	}

	if word.WordDate.IsZero() {
		return errors.New("word_date is required")
	}

	return nil
}
