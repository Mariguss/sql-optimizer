package recommendation

import (
	"fmt"
	//"strings"
)

// Engine генерирует рекомендации на основе анализа
type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

// GenerateRecommendations генерирует рекомендации на основе проблемных операций
func (e *Engine) GenerateRecommendations(result *AnalysisResult) []string {
	var recommendations []string

	// Добавляем общие рекомендации
	recommendations = append(recommendations, e.generateGeneralRecommendations(result)...)

	// Добавляем специфические рекомендации для проблемных операций
	for _, problem := range result.ProblematicOperations {
		recommendations = append(recommendations, e.generateSpecificRecommendations(problem)...)
	}

	// Убираем дубликаты
	return e.removeDuplicates(recommendations)
}

// generateGeneralRecommendations общие рекомендации
func (e *Engine) generateGeneralRecommendations(result *AnalysisResult) []string {
	var recs []string

	if result.TotalCost > 10000 {
		recs = append(recs, "Общая стоимость запроса очень высока. Рассмотрите рефакторинг запроса или добавление индексов.")
	}

	if result.TotalActualTime != nil && *result.TotalActualTime > 1000 {
		recs = append(recs, "Общее время выполнения превышает 1 секунду. Оптимизация необходима.")
	}

	return recs
}

// generateSpecificRecommendations специфические рекомендации для типов операций
func (e *Engine) generateSpecificRecommendations(problem ProblematicOperation) []string {
	var recs []string

	switch problem.NodeType {
	case "Seq Scan":
		recs = append(recs, e.handleSeqScan(problem))
	case "Sort":
		recs = append(recs, e.handleSort(problem))
	case "Hash Join":
		recs = append(recs, e.handleHashJoin(problem))
	case "Nested Loop":
		recs = append(recs, e.handleNestedLoop(problem))
	default:
		recs = append(recs, fmt.Sprintf("Операция %s имеет высокую стоимость. Рассмотрите оптимизацию.", problem.NodeType))
	}

	return recs
}

func (e *Engine) handleSeqScan(problem ProblematicOperation) string {
	return "Sequential Scan обнаружен. Добавьте индексы на поля, используемые в условиях фильтрации."
}

func (e *Engine) handleSort(problem ProblematicOperation) string {
	return "Обнаружена операция сортировки. Используйте индексы для предварительной сортировки данных."
}

func (e *Engine) handleHashJoin(problem ProblematicOperation) string {
	return "Hash Join обнаружен. Убедитесь, что обе таблицы имеют индексы на полях соединения."
}

func (e *Engine) handleNestedLoop(problem ProblematicOperation) string {
	return "Nested Loop обнаружен. Рассмотрите изменение условий соединения или добавление индексов."
}

// removeDuplicates убирает дублирующиеся рекомендации
func (e *Engine) removeDuplicates(recommendations []string) []string {
	seen := make(map[string]bool)
	var unique []string

	for _, rec := range recommendations {
		if !seen[rec] {
			seen[rec] = true
			unique = append(unique, rec)
		}
	}

	return unique
}

// AnalysisResult представляет результат анализа запроса
type AnalysisResult struct {
	TotalCost             float64
	TotalActualTime       *float64
	ProblematicOperations []ProblematicOperation
	Recommendations       []string
	Warnings              []string
}

// ProblematicOperation представляет проблемную операцию
type ProblematicOperation struct {
	NodeType       string
	Cost           float64
	ActualTime     *float64
	Description    string
	Recommendation string
	Severity       string
}
