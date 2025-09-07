package analyzer

import (
	"encoding/json"
	"fmt"
)

// Структура для полного ответа EXPLAIN
type ExplainResult struct {
	Plan         PlanNode                 `json:"Plan"`
	Planning     map[string]interface{}   `json:"Planning"`
	PlanningTime float64                  `json:"Planning Time"`
	Triggers     []interface{}            `json:"Triggers"`
	ExecutionTime float64                 `json:"Execution Time"`
}

// PlanNode представляет узел плана выполнения PostgreSQL
type PlanNode struct {
	NodeType          string    `json:"Node Type"`
	RelationName      string    `json:"Relation Name,omitempty"`
	Alias             string    `json:"Alias,omitempty"`
	StartupCost       float64   `json:"Startup Cost"`
	TotalCost         float64   `json:"Total Cost"`
	PlanRows          int       `json:"Plan Rows"`
	PlanWidth         int       `json:"Plan Width,omitempty"`
	ActualStartupTime *float64  `json:"Actual Startup Time,omitempty"`
	ActualTotalTime   *float64  `json:"Actual Total Time,omitempty"`
	ActualRows        *int      `json:"Actual Rows,omitempty"`
	ActualLoops       *int      `json:"Actual Loops,omitempty"`
	Filter            string    `json:"Filter,omitempty"`
	JoinType          string    `json:"Join Type,omitempty"`
	IndexName         string    `json:"Index Name,omitempty"`
	HashCondition     string    `json:"Hash Condition,omitempty"`
	Plans             []PlanNode `json:"Plans,omitempty"`
}

// ParseExplainJSON парсит JSON вывод EXPLAIN
func ParseExplainJSON(planJSON string) ([]PlanNode, error) {
	fmt.Printf("🔍 Парсим JSON...\n")
	
	// Пробуем парсить как массив ExplainResult
	var explainResults []ExplainResult
	if err := json.Unmarshal([]byte(planJSON), &explainResults); err == nil {
		fmt.Printf("✅ Распаршено как массив ExplainResult: %d элементов\n", len(explainResults))
		var planNodes []PlanNode
		for _, result := range explainResults {
			planNodes = append(planNodes, result.Plan)
		}
		return planNodes, nil
	}

	// Пробуем парсить как одиночный ExplainResult
	var singleExplainResult ExplainResult
	if err := json.Unmarshal([]byte(planJSON), &singleExplainResult); err == nil {
		fmt.Printf("✅ Распаршено как одиночный ExplainResult\n")
		return []PlanNode{singleExplainResult.Plan}, nil
	}

	// Пробуем парсить как массив PlanNode (старый формат)
	var planNodes []PlanNode
	if err := json.Unmarshal([]byte(planJSON), &planNodes); err == nil {
		fmt.Printf("✅ Распаршено как массив PlanNode: %d элементов\n", len(planNodes))
		return planNodes, nil
	}

	// Пробуем парсить как одиночный PlanNode
	var singlePlan PlanNode
	if err := json.Unmarshal([]byte(planJSON), &singlePlan); err == nil {
		fmt.Printf("✅ Распаршено как одиночный PlanNode\n")
		return []PlanNode{singlePlan}, nil
	}

	return nil, fmt.Errorf("не удалось распарсить JSON: %v", planJSON)
}