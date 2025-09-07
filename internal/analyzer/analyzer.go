package analyzer

import (
	"fmt"
)

// AnalyzePlan анализирует план выполнения из JSON
func AnalyzePlan(planJSON string) (*AnalysisResult, error) {
	planNodes, err := ParseExplainJSON(planJSON)
	if err != nil {
		return nil, err
	}

	result := &AnalysisResult{
		ProblematicOperations: []ProblematicOperation{},
		Recommendations:       []string{},
		Warnings:              []string{},
	}

	// Анализируем все узлы плана
	for _, plan := range planNodes {
		analyzeNode(plan, result)
	}

	fmt.Printf("📊 Результат анализа: TotalCost=%.2f, Problems=%d\n", 
		result.TotalCost, len(result.ProblematicOperations))

	return result, nil
}

// analyzeNode рекурсивно анализирует узлы плана
func analyzeNode(node PlanNode, result *AnalysisResult) {
	// Анализируем текущий узел
	analyzeSingleNode(node, result)

	// Рекурсивно анализируем дочерние узлы
	for _, child := range node.Plans {
		analyzeNode(child, result)
	}
}

// analyzeSingleNode анализирует одиночный узел
func analyzeSingleNode(node PlanNode, result *AnalysisResult) {
	result.TotalCost += node.TotalCost

	if node.ActualTotalTime != nil {
		if result.TotalActualTime == nil {
			result.TotalActualTime = new(float64)
		}
		*result.TotalActualTime += *node.ActualTotalTime
	}

	// Проверяем проблемные операции
	switch node.NodeType {
	case "Seq Scan":
		if node.TotalCost > 1.0 { // Понизим порог для теста
			problem := ProblematicOperation{
				NodeType:      "Seq Scan",
				Cost:          node.TotalCost,
				ActualTime:    node.ActualTotalTime,
				Description:   fmt.Sprintf("Sequential Scan на таблице %s", node.RelationName),
				Recommendation: "Добавить индекс на используемые в WHERE поля",
				Severity:      "high",
			}
			result.ProblematicOperations = append(result.ProblematicOperations, problem)
			result.Recommendations = append(result.Recommendations, 
				fmt.Sprintf("Создать индекс для таблицы %s", node.RelationName))
		}

	case "Sort":
		if node.TotalCost > 0.5 {
			problem := ProblematicOperation{
				NodeType:      "Sort",
				Cost:          node.TotalCost,
				ActualTime:    node.ActualTotalTime,
				Description:   "Операция сортировки",
				Recommendation: "Использовать индексы для предварительной сортировки",
				Severity:      "medium",
			}
			result.ProblematicOperations = append(result.ProblematicOperations, problem)
		}

	case "Hash Join", "Nested Loop":
		if node.TotalCost > 2.0 {
			problem := ProblematicOperation{
				NodeType:      node.NodeType,
				Cost:          node.TotalCost,
				ActualTime:    node.ActualTotalTime,
				Description:   fmt.Sprintf("Операция соединения %s", node.NodeType),
				Recommendation: "Проверить индексы на полях соединения",
				Severity:      "medium",
			}
			result.ProblematicOperations = append(result.ProblematicOperations, problem)
		}
	}

	// Дополнительные проверки
	if node.PlanRows > 1000 && node.ActualRows != nil && *node.ActualRows < node.PlanRows/10 {
		result.Warnings = append(result.Warnings, 
			fmt.Sprintf("Плохая оценка строк: планировалось %d, фактически %d", 
				node.PlanRows, *node.ActualRows))
	}
}

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