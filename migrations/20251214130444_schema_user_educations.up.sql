CREATE TABLE user_educations (
                                 id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '学歴ID',
                                 user_id BIGINT UNSIGNED NOT NULL COMMENT 'ユーザーID（users.id）',
                                 institution_name VARCHAR(255) NOT NULL COMMENT '学校名',
                                 faculty_name VARCHAR(255) NULL COMMENT '学部',
                                 department_name VARCHAR(255) NULL COMMENT '学科・専攻',
                                 degree_type_id BIGINT UNSIGNED NULL COMMENT '学位種別ID（degree_types.id）',
                                 education_status_id BIGINT UNSIGNED NOT NULL COMMENT '学歴状態ID（education_statuses.id）',
                                 event_date DATE NOT NULL COMMENT '年月',
                                 description TEXT NULL COMMENT '補足',
                                 sort_order INT NOT NULL DEFAULT 0 COMMENT '表示順',
                                 is_public BOOLEAN NOT NULL DEFAULT TRUE COMMENT '公開可否',
                                 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
                                 updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

                                 PRIMARY KEY (id),

    -- 検索・一覧取得用
                                 KEY idx_user_educations_user_sort (user_id, sort_order),
                                 KEY idx_user_educations_user (user_id),

    -- 外部キー
                                 CONSTRAINT fk_user_educations_user
                                     FOREIGN KEY (user_id) REFERENCES users(id)
                                         ON DELETE CASCADE,

                                 CONSTRAINT fk_user_educations_degree_type
                                     FOREIGN KEY (degree_type_id) REFERENCES degree_types(id),

                                 CONSTRAINT fk_user_educations_education_status
                                     FOREIGN KEY (education_status_id) REFERENCES education_statuses(id)

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='ユーザー学歴';
