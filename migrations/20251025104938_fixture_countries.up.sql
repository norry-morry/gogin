INSERT INTO countries (code, is_supported, sort_order, phone_code)
VALUES
    ('JP', 1, 1, '+81'),     -- 日本
    ('US', 1, 2, '+1'),      -- アメリカ
    ('FR', 0, 3, '+33'),      -- フランス (+33)
    ('ES', 0, 4, '+34'),      -- スペイン (+34)
    ('DE', 0, 5, '+49'),      -- ドイツ (+49)
    ('CN', 0, 6, '+86'),      -- 中国 (+86)
    ('KR', 0, 7, '+82')      -- 韓国 (+82)
ON DUPLICATE KEY UPDATE
    sort_order = VALUES(sort_order),
    is_supported = VALUES(is_supported);
