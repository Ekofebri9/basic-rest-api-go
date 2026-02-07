package handlers

import (
	"basic-rest-api-go/internal/dto"
	"basic-rest-api-go/internal/services"
	"encoding/json"
	"net/http"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GenerateReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var req dto.ReportRequest

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// default range in today's date
	if startDate == "" || endDate == "" {
		now := time.Now()
		startDate = now.Format("2006-01-02")
		endDate = now.AddDate(0, 0, 1).Format("2006-01-02")
	}

	req.StartDate = startDate
	req.EndDate = endDate

	transaction, err := h.service.GenerateSalesReport(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(transaction)
}
