package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http" // Add this line
	"sql-optimizer/internal/analyzer"
	"sql-optimizer/internal/postgres" // Add this line
	"time"

	_ "github.com/lib/pq"
)

// DBConfig хранит параметры подключения к БД
type DBConfig struct {
	Host     string `json:"db_host"`
	Port     string `json:"db_port"`
	User     string `json:"db_user"`
	Password string `json:"db_password"`
	DBName   string `json:"db_name"`
}

// AnalyzeRequest представляет структуру запроса для анализа
type AnalyzeRequest struct {
	DBConfig
	Query string `json:"query"`
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// getDB открывает и настраивает соединение с БД.
// Важно: эта функция не закрывает соединение, вызывающий код должен это сделать.
func getDB(config DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вызове sql.Open: %w", err)
	}

	// Настройки пула соединений
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	return db, nil
}

// ConnectDB теперь просто проверяет, что соединение может быть установлено.
func (h *Handler) ConnectDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var config DBConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка парсинга JSON: %v", err), http.StatusBadRequest)
		return
	}

	db, err := getDB(config)
	if err != nil {
		http.Error(w, "Ошибка подготовки к подключению к БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close() // Закрываем соединение сразу после проверки

	// Проверяем соединение с помощью Ping
	if err = db.Ping(); err != nil {
		http.Error(w, "Ошибка ping БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Успешное подключение к базе данных!"})
}

// AnalyzeQuery открывает соединение, анализирует запрос и закрывает соединение.
func (h *Handler) AnalyzeQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка парсинга JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Создаем клиент к БД, который будет закрыт в конце
	pgClient, err := postgres.NewClient(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		req.DBConfig.Host, req.DBConfig.Port, req.DBConfig.User, req.DBConfig.Password, req.DBConfig.DBName))
	if err != nil {
		http.Error(w, "Ошибка подключения к БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer pgClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Получаем JSON план выполнения
	planJSON, err := pgClient.GetExplainPlan(ctx, req.Query)
	if err != nil {
		http.Error(w, "Ошибка получения плана EXPLAIN: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Анализируем план и получаем результат
	analysisResult, err := analyzer.AnalyzePlan(planJSON)
	if err != nil {
		http.Error(w, "Ошибка анализа плана: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем полный результат анализа в формате JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analysisResult); err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
