package util

import (
	"encoding/json"
	"net/http"
)

func SendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")

	if code == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	}

	w.WriteHeader(code)

	response := map[string]any{
		"error":   http.StatusText(code),
		"message": message,
		"code":    code,
	}

	json.NewEncoder(w).Encode(response)
}

func SendSuccess(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := map[string]any{
		"success": true,
		"data":    data,
	}

	json.NewEncoder(w).Encode(response)
}

func SendAuthSuccess(w http.ResponseWriter, userID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]any{
		"success": true,
		"user_id": userID,
	}

	json.NewEncoder(w).Encode(response)
}
