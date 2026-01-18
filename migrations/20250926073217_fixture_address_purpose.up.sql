INSERT INTO address_purposes
    (id, code, sort_order, is_active)
VALUES
    (1,'home', 10,    1),
    (2, 'contact', 20,   1),
    (3, 'office', 30,  1),
    (4, 'shipping', 40, 1),
    (5, 'billing', 50, 1),
    (6, 'other', 90,   1)
ON DUPLICATE KEY UPDATE
    sort_order = VALUES(sort_order),
                     is_active  = VALUES(is_active);

