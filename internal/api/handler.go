package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sql-optimizer/internal/analyzer"
	"sql-optimizer/internal/postgres"
)

type Handler struct {
	pgClient *postgres.Client
}

func NewHandler(pgClient *postgres.Client) *Handler {
	return &Handler{pgClient: pgClient}
}

func (h *Handler) AnalyzeQuery(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query string `json:"query"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("📨 Получен запрос: %s\n", req.Query)

	// Получаем план выполнения
	planJSON, err := h.pgClient.GetExplainPlan(r.Context(), req.Query)
	if err != nil {
		fmt.Printf("❌ Ошибка получения плана: %v\n", err)
		sendError(w, "Failed to get explain plan: "+err.Error())
		return
	}

	fmt.Printf("📊 Получен план: %s\n", planJSON)

	// Анализируем план
	result, err := analyzer.AnalyzePlan(planJSON)
	if err != nil {
		fmt.Printf("❌ Ошибка анализа: %v\n", err)
		sendError(w, "Failed to analyze plan: "+err.Error())
		return
	}

	fmt.Printf("✅ Результат анализа: %+v\n", result)

	// Отправляем результат
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

func sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}