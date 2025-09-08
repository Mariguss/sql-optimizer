package analyzer

import (
    "reflect"
    "testing"
)

// type AnalysisResult struct {
//     TotalCost            float64
//     TotalActualTime      *float64
//     ProblematicOperations []ProblematicOperation
//     Recommendations      []string
//     Warnings             []string
// }

// type ProblematicOperation struct {
//     NodeType       string
//     Cost           float64
//     ActualTime     *float64
//     Description    string
//     Recommendation string
//     Severity       string
// }

func TestAnalyzePlan(t *testing.T) {
    testCases := []struct {
        name        string
        planJSON    string
        expected    *AnalysisResult
        expectError bool
    }{
        {
            name: "Успешный анализ Seq Scan",
            planJSON: `[
                {
                    "Plan": {
                        "Node Type": "Seq Scan",
                        "Relation Name": "users",
                        "Total Cost": 150.5,
                        "Plan Rows": 10000,
                        "Actual Total Time": 25.3,
                        "Actual Rows": 10000
                    }
                }
            ]`,
            expected: &AnalysisResult{
                TotalCost:       150.5,
                TotalActualTime: float64Ptr(25.3),
                ProblematicOperations: []ProblematicOperation{
                    {
                        NodeType:      "Seq Scan",
                        Cost:          150.5,
                        ActualTime:    float64Ptr(25.3),
                        Description:   "Sequential Scan на таблице users",
                        Recommendation: "Добавить индекс на используемые в WHERE поля",
                        Severity:      "high",
                    },
                },
                Recommendations: []string{"Создать индекс для таблицы users"},
                Warnings:        []string{},
            },
            expectError: false,
        },
        {
            name: "Анализ с операцией Sort",
            planJSON: `[
                {
                    "Plan": {
                        "Node Type": "Sort",
                        "Total Cost": 10.0,
                        "Plans": []
                    }
                }
            ]`,
            expected: &AnalysisResult{
                TotalCost: 10.0,
                ProblematicOperations: []ProblematicOperation{
                    {
                        NodeType:      "Sort",
                        Cost:          10.0,
                        Description:   "Операция сортировки",
                        Recommendation: "Использовать индексы для предварительной сортировки",
                        Severity:      "medium",
                    },
                },
                Recommendations: []string{},
                Warnings:        []string{},
            },
            expectError: false,
        },
        {
            name: "Невалидный JSON",
            planJSON:    `{"invalid json`,
            expected:    nil,
            expectError: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := AnalyzePlan(tc.planJSON)

            if tc.expectError {
                if err == nil {
                    t.Errorf("Ожидалась ошибка, но получили nil")
                }
                return
            }

            if err != nil {
                t.Fatalf("Неожиданная ошибка: %v", err)
            }

            if !reflect.DeepEqual(result, tc.expected) {
                t.Errorf("Результат не соответствует ожидаемому.\nОжидали: %+v\nПолучили: %+v", tc.expected, result)
            }
        })
    }
}

func float64Ptr(v float64) *float64 {
    return &v
}