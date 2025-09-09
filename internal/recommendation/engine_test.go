package recommendation

import (
	"reflect"
	"sort"
	"testing"
)

func TestGenerateRecommendations(t *testing.T) {
	engine := NewEngine()

	// Вспомогательная функция для создания указателя на float64
	float64Ptr := func(v float64) *float64 {
		return &v
	}

	testCases := []struct {
		name           string
		analysisResult *AnalysisResult
		expectedRecs   []string
	}{
		{
			name: "Общая рекомендация по высокой стоимости",
			analysisResult: &AnalysisResult{
				TotalCost: 15000,
			},
			expectedRecs: []string{
				"Общая стоимость запроса очень высока. Рассмотрите рефакторинг запроса или добавление индексов.",
			},
		},
		{
			name: "Общая рекомендация по времени выполнения",
			analysisResult: &AnalysisResult{
				TotalActualTime: float64Ptr(1200),
			},
			expectedRecs: []string{
				"Общее время выполнения превышает 1 секунду. Оптимизация необходима.",
			},
		},
		{
			name: "Специфическая рекомендация для Seq Scan",
			analysisResult: &AnalysisResult{
				ProblematicOperations: []ProblematicOperation{
					{NodeType: "Seq Scan"},
				},
			},
			expectedRecs: []string{
				"Sequential Scan обнаружен. Добавьте индексы на поля, используемые в условиях фильтрации.",
			},
		},
		{
			name: "Специфическая рекомендация для Sort",
			analysisResult: &AnalysisResult{
				ProblematicOperations: []ProblematicOperation{
					{NodeType: "Sort"},
				},
			},
			expectedRecs: []string{
				"Обнаружена операция сортировки. Используйте индексы для предварительной сортировки данных.",
			},
		},
		{
			name: "Комбинация общей и специфической рекомендаций",
			analysisResult: &AnalysisResult{
				TotalCost: 20000,
				ProblematicOperations: []ProblematicOperation{
					{NodeType: "Nested Loop"},
				},
			},
			expectedRecs: []string{
				"Nested Loop обнаружен. Рассмотрите изменение условий соединения или добавление индексов.",
				"Общая стоимость запроса очень высока. Рассмотрите рефакторинг запроса или добавление индексов.",
			},
		},
		{
			name: "Удаление дубликатов",
			analysisResult: &AnalysisResult{
				ProblematicOperations: []ProblematicOperation{
					{NodeType: "Seq Scan"},
					{NodeType: "Seq Scan"},
				},
			},
			expectedRecs: []string{
				"Sequential Scan обнаружен. Добавьте индексы на поля, используемые в условиях фильтрации.",
			},
		},
		{
			name: "Неизвестный тип операции",
			analysisResult: &AnalysisResult{
				ProblematicOperations: []ProblematicOperation{
					{NodeType: "Unknown Operation"},
				},
			},
			expectedRecs: []string{
				"Операция Unknown Operation имеет высокую стоимость. Рассмотрите оптимизацию.",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recs := engine.GenerateRecommendations(tc.analysisResult)

			// Сортируем оба слайса, так как порядок не гарантирован
			sort.Strings(recs)
			sort.Strings(tc.expectedRecs)

			if !reflect.DeepEqual(recs, tc.expectedRecs) {
				t.Errorf("Рекомендации не соответствуют ожидаемым.\nОжидали: %v\nПолучили: %v", tc.expectedRecs, recs)
			}
		})
	}
}
