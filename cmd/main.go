package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"sql-optimizer/internal/api"
	"sql-optimizer/internal/postgres"
)

func main() {
	fmt.Println("SQL Optimizer - Анализатор запросов PostgreSQL")
	fmt.Println("==============================================")

	// Получаем параметры подключения
	connStr := getConnectionInfo()

	fmt.Println("\nПодключаемся к PostgreSQL...")
	fmt.Println("Строка подключения:", connStr)

	// Подключаемся к PostgreSQL
	pgClient, err := postgres.NewClient(connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}
	defer pgClient.Close()

	fmt.Println("Успешно подключено к PostgreSQL!")

	// Создаем обработчик
	handler := api.NewHandler(pgClient)

	// Настраиваем маршруты
	http.HandleFunc("/api/analyze", handler.AnalyzeQuery)

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", http.StripPrefix("/", fs))
	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println(" Откройте этот адрес в браузере")
	fmt.Println("\nДля остановки нажмите Ctrl+C")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getConnectionInfo() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n Введите параметры подключения к PostgreSQL:")
	fmt.Println("   (оставьте пустым для значений по умолчанию или используйте переменные окружения)")
	fmt.Println()

	// Хост
	host := os.Getenv("DB_HOST")
	if host == "" {
		fmt.Print(" Хост сервера [localhost]: ")
		host, _ = reader.ReadString('\n')
		host = strings.TrimSpace(host)
		if host == "" {
			host = "localhost"
		}
	} else {
		fmt.Printf(" Хост сервера (из окружения): %s\n", host)
	}

	// Порт
	port := os.Getenv("DB_PORT")
	if port == "" {
		fmt.Print("Порт [5432]: ")
		port, _ = reader.ReadString('\n')
		port = strings.TrimSpace(port)
		if port == "" {
			port = "5432"
		}
	} else {
		fmt.Printf("Порт (из окружения): %s\n", port)
	}

	// Имя пользователя
	user := os.Getenv("DB_USER")
	if user == "" {
		fmt.Print("Имя пользователя [postgres]: ")
		user, _ = reader.ReadString('\n')
		user = strings.TrimSpace(user)
		if user == "" {
			user = "postgres"
		}
	} else {
		fmt.Printf("Имя пользователя (из окружения): %s\n", user)
	}

	// Пароль
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		fmt.Print("Пароль: ")
		password, _ = reader.ReadString('\n')
		password = strings.TrimSpace(password)
	} else {
		fmt.Printf("Пароль (из окружения): %s\n", strings.Repeat("*", len(password)))
	}

	// Имя базы данных
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		fmt.Print("Имя базы данных [postgres]: ")
		dbname, _ = reader.ReadString('\n')
		dbname = strings.TrimSpace(dbname)
		if dbname == "" {
			dbname = "postgres"
		}
	} else {
		fmt.Printf("Имя базы данных (из окружения): %s\n", dbname)
	}

	// SSL mode
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		fmt.Print("SSL mode [disable]: ")
		sslmode, _ = reader.ReadString('\n')
		sslmode = strings.TrimSpace(sslmode)
		if sslmode == "" {
			sslmode = "disable"
		}
	} else {
		fmt.Printf("SSL mode (из окружения): %s\n", sslmode)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}
