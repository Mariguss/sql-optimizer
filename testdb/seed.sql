-- Вставка данных в таблицу users
INSERT INTO users (name, email) VALUES
('Alex', 'alex@example.com'),
('Maria', 'maria@example.com'),
('Ivan', 'ivan@example.com'),
('Anna', 'anna@example.com'),
('Sergey', 'sergey@example.com');

-- Вставка данных в таблицу orders
INSERT INTO orders (user_id, amount, order_date, status) VALUES
(1, 150.00, '2025-09-01', 'completed'),
(1, 25.50, '2025-09-02', 'pending'),
(2, 200.00, '2025-09-03', 'completed'),
(3, 75.20, '2025-09-04', 'shipped'),
(4, 120.00, '2025-09-05', 'completed'),
(5, 50.00, '2025-09-06', 'pending');

-- Вставка большого объема данных для тестирования
DO $$
BEGIN
    FOR i IN 1..10000 LOOP
        INSERT INTO users (name, email)
        VALUES (
            format('User%s', i),
            format('user%s@test.com', i)
        );
    END LOOP;
END;
$$;