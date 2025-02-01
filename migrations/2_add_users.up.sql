INSERT INTO users (name) VALUES
('Sheldon'),
('Leonard'),
('Penny'),
('Howard'),
('Rajesh');


-- SELECT
--     t.id AS transaction_id,
--     CASE
--         WHEN t.amount < 100 THEN TO_CHAR(t.amount, 'FM9999999999999999999990.00')  -- Если меньше 100, добавляем ".00"
--         ELSE TO_CHAR(t.amount / 100.0, 'FM9999999999999999999990.00')  -- Иначе делим на 100
--     END AS amount,
--     t.date_operation,
--     sender.name AS sender_name,
--     receiver.name AS receiver_name
-- FROM
--     transactions t
-- JOIN
--     users sender ON t.sender_id = sender.id
-- JOIN
--     users receiver ON t.receiver_id = receiver.id
-- WHERE
--     t.sender_id = 1 OR t.receiver_id = 1
-- ORDER BY
--     t.date_operation DESC
-- LIMIT 10;