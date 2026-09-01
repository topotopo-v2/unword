package word

import (
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) GetToday(w http.ResponseWriter, r *http.Request) {
	timezone := r.URL.Query().Get("timezone")

	if timezone == "" {
		http.Error(
			w,
			`{"error":"timezone is required"}`,
			http.StatusBadRequest,
		)
		return
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		http.Error(
			w,
			`{"error":"invalid timezone"}`,
			http.StatusBadRequest,
		)
		return
	}

	today := time.Now().In(location)

	word, err := h.repository.GetByDate(
		r.Context(),
		today,
	)

	if err != nil {
		http.Error(
			w,
			`{"error":"word not found"}`,
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(word)
}
