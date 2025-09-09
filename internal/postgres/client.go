package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Client struct {
	db *sql.DB
}

func NewClient(connectionString string) (*Client, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping БД: %v", err)
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

// GetExplainPlan получает план выполнения в формате JSON

func (c *Client) GetExplainPlan(ctx context.Context, query string) (string, error) {
	var planJSON string

	explainQuery := fmt.Sprintf("EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) %s", query)

	fmt.Printf("Выполняем: %s\n", explainQuery)

	err := c.db.QueryRowContext(ctx, explainQuery).Scan(&planJSON)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения EXPLAIN: %v", err)
	}

	fmt.Printf("Получен JSON плана: %s\n", planJSON)
	return planJSON, nil
}
