# SQL Optimizer - Анализатор запросов PostgreSQL

Инструмент для анализа и оптимизации SQL запросов к PostgreSQL

## Возможности

- Парсинг планов выполнения EXPLAIN
- Выявление проблемных операций (Seq Scan, Sort, Hash Join)
- Рекомендации по оптимизации запросов
- Веб-интерфейс для удобной работы

## Установка и запуск

### Требования
- Go 1.21+
- PostgreSQL 10+

### 1. Клонирование репозитория
```bash
git clone https://github.com/your-username/sql-optimizer.git
cd sql-optimizer
