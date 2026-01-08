INSERT INTO degree_types (id, code, sort_order, is_active) VALUES
                                                           (1, 'high_school',    1, TRUE),
                                                           (2, 'vocational',     2, TRUE),
                                                           (3, 'junior_college', 3, TRUE),
                                                           (4, 'bachelor',       4, TRUE),
                                                           (5, 'master',         5, TRUE),
                                                           (6, 'doctor',         6, TRUE),
                                                           (7, 'other',          7, TRUE)
    ON DUPLICATE KEY UPDATE
                         sort_order = VALUES(sort_order),
                         is_active  = VALUES(is_active);
