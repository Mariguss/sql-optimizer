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
	fmt.Println("🐘 SQL Optimizer - Анализатор запросов PostgreSQL")
	fmt.Println("==============================================")
	
	// Получаем параметры подключения
	connStr := getConnectionInfo()
	
	fmt.Println("\n🔗 Подключаемся к PostgreSQL...")
	fmt.Println("Строка подключения:", connStr)
	
	// Подключаемся к PostgreSQL
	pgClient, err := postgres.NewClient(connStr)
	if err != nil {
		log.Fatal("❌ Ошибка подключения к PostgreSQL:", err)
	}
	defer pgClient.Close()

	fmt.Println("✅ Успешно подключено к PostgreSQL!")

	// Создаем обработчик
	handler := api.NewHandler(pgClient)

	// Настраиваем маршруты
	http.HandleFunc("/api/analyze", handler.AnalyzeQuery)
	
	// Простая HTML страница прямо в коде
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>SQL Optimizer</title>
				<meta charset="UTF-8">
				<style>
					body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
					.container { max-width: 800px; margin: 0 auto; background: white; padding: 20px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
					textarea { width: 100%; height: 100px; padding: 10px; border: 1px solid #ddd; border-radius: 5px; }
					button { padding: 10px 20px; background: #007bff; color: white; border: none; border-radius: 5px; cursor: pointer; }
					.result { margin-top: 20px; padding: 15px; background: #f8f9fa; border-radius: 5px; }
				</style>
			</head>
			<body>
				<div class="container">
					<h1>🔍 SQL Optimizer</h1>
					<p>Введите SQL запрос для анализа:</p>
					<textarea id="query" placeholder="SELECT * FROM users WHERE age > 25"></textarea>
					<br>
					<button onclick="analyze()">Анализировать запрос</button>
					<div id="result" class="result">Результат появится здесь...</div>
				</div>
				
				<script>
					function analyze() {
						const query = document.getElementById('query').value;
						const resultDiv = document.getElementById('result');
						
						if (!query) {
							resultDiv.innerHTML = '❌ Введите SQL запрос';
							return;
						}
						
						resultDiv.innerHTML = '⏳ Анализируем...';
						
						fetch('/api/analyze', {
							method: 'POST',
							headers: {'Content-Type': 'application/json'},
							body: JSON.stringify({query: query})
						})
						.then(response => response.json())
						.then(data => {
							resultDiv.innerHTML = '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
						})
						.catch(error => {
							resultDiv.innerHTML = '❌ Ошибка: ' + error;
						});
					}
				</script>
			</body>
			</html>
		`))
	})

	fmt.Println("🌐 Сервер запущен на http://localhost:8080")
	fmt.Println("📋 Откройте этот адрес в браузере")
	fmt.Println("\n🛑 Для остановки нажмите Ctrl+C")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getConnectionInfo() string {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("\n📝 Введите параметры подключения к PostgreSQL:")
	fmt.Println("   (оставьте пустым для значений по умолчанию)")
	fmt.Println()

	fmt.Print("🏠 Хост сервера [localhost]: ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)
	if host == "" { host = "localhost" }

	fmt.Print("🚪 Порт [5432]: ")
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)
	if port == "" { port = "5432" }

	fmt.Print("👤 Имя пользователя [postgres]: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)
	if user == "" { user = "postgres" }

	fmt.Print("🔑 Пароль: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("💾 Имя базы данных [postgres]: ")
	dbname, _ := reader.ReadString('\n')
	dbname = strings.TrimSpace(dbname)
	if dbname == "" { dbname = "postgres" }

	fmt.Print("🔒 SSL mode [disable]: ")
	sslmode, _ := reader.ReadString('\n')
	sslmode = strings.TrimSpace(sslmode)
	if sslmode == "" { sslmode = "disable" }

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}