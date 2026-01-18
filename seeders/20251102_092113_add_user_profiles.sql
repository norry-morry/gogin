INSERT INTO user_profiles
(user_id, family_name, given_name, family_name_kana, given_name_kana, birth_date, gender_id, initial)
VALUES
    (1, '森', '紀洋', 'モリ', 'ノリヒロ', '1975-06-17', 1, 'n.m')
    ON DUPLICATE KEY UPDATE
                         family_name=VALUES(family_name),
                         given_name=VALUES(given_name),
                         family_name_kana=VALUES(family_name_kana),
                         given_name_kana=VALUES(given_name_kana),
                         birth_date=VALUES(birth_date),
                         gender_id=VALUES(gender_id),
                         initial=VALUES(initial),
                         updated_at=CURRENT_TIMESTAMP(6);
