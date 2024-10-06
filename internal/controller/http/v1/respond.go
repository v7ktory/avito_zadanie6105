package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, fmt.Sprintf("JSON encode error: %s", err.Error()), http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"reason": message})
}

func respondWithValidationErrors(w http.ResponseWriter, problems map[string]string) {
	for _, problem := range problems {
		respondWithError(w, http.StatusBadRequest, problem)
	}
}
