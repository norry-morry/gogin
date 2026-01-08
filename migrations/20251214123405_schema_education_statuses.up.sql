CREATE TABLE education_statuses (
                                    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
                                    code VARCHAR(32) NOT NULL COMMENT '内部コード',
                                    sort_order INT NOT NULL DEFAULT 0 COMMENT '並び順',
                                    is_active BOOLEAN NOT NULL DEFAULT TRUE COMMENT '有効フラグ',
                                    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
                                    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
                                    PRIMARY KEY (id),
                                    UNIQUE KEY uk_education_statuses_code (code),
                                    KEY idx_education_statuses_active_sort (is_active, sort_order)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='学歴状態マスタ';
