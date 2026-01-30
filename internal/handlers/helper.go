package handlers

import (
	"basic-rest-api-go/internal/models"
	"encoding/json"
	"net/http"
)

func parseBody(w http.ResponseWriter, r *http.Request, v *models.Product) {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
}
