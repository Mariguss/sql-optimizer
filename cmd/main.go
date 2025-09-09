package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"sql-optimizer/internal/api"
)

func main() {
	fmt.Println("SQL Optimizer - Web-сервер запущен")
	fmt.Println("===================================")

	handler := api.NewHandler()

	http.HandleFunc("/api/connect", handler.ConnectDB)
	http.HandleFunc("/api/analyze", handler.AnalyzeQuery)

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Сервер доступен по адресу: http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
