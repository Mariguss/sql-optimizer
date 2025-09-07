package types

// AnalysisResult представляет результат анализа запроса
type AnalysisResult struct {
	TotalCost             float64                 `json:"total_cost"`
	TotalActualTime       *float64                `json:"total_actual_time,omitempty"`
	ProblematicOperations []ProblematicOperation  `json:"problematic_operations"`
	Recommendations       []string                `json:"recommendations"`
	Warnings              []string                `json:"warnings"`
}

// ProblematicOperation представляет проблемную операцию
type ProblematicOperation struct {
	NodeType      string   `json:"node_type"`
	Cost          float64  `json:"cost"`
	ActualTime    *float64 `json:"actual_time,omitempty"`
	Description   string   `json:"description"`
	Recommendation string   `json:"recommendation"`
	Severity      string   `json:"severity"` // "high", "medium", "low"
}

// ExplainPlanRequest запрос на анализ
type ExplainPlanRequest struct {
	Query string `json:"query"`
}

// ExplainPlanResponse ответ с анализом
type ExplainPlanResponse struct {
	Success bool            `json:"success"`
	Data    *AnalysisResult `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}