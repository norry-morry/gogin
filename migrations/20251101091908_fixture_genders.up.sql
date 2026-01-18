INSERT INTO genders
    (id, code, sort_order)
VALUES
    (1, 'male',        1),
    (2, 'female',      2),
    (3, 'other',       3),
    (4, 'unspecified', 9)
ON DUPLICATE KEY UPDATE
                     sort_order = VALUES(sort_order);
