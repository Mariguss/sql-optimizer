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
```
### 2. Запуск docker compose
```bash
docker-compose up --build
```
### 3. Открыть на localhost
Либо уже развернутое http://185.159.111.69:8081/

### 4.Подключение
Подключитесь к базе данных по ip облачной бд, порту, пользователю, паролю, имени базы данных

### 5.Ссылка на яндекс диск
https://disk.yandex.ru/d/05B6JRSrTCFShA

