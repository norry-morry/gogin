INSERT INTO education_statuses
    (id, code, sort_order)
VALUES
    (1, 'entrance', 1),
    (2, 'enrolled', 2),
    (3, 'leave_of_absence', 3), -- 在学の次に休学
    (4, 'graduated', 4),
    (5, 'completed', 5),
    (6, 'graduation_prospect', 6),
    (7, 'withdrawn', 7),
    (8, 'expelled', 8)
ON DUPLICATE KEY UPDATE
                     sort_order = VALUES(sort_order),
                     is_active  = VALUES(is_active);
