-- 例: 20251101_200001_genders.up.sql
CREATE TABLE genders (
                         id          TINYINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '性別マスタID(PK)',
                         code        VARCHAR(32) NOT NULL UNIQUE COMMENT 'male/female/other/unspecified',
                         sort_order  TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '並び順',
                         is_active   BOOLEAN NOT NULL DEFAULT TRUE COMMENT '有効化フラグ',
                         created_at  DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
                         updated_at  DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時'
)
    ENGINE=InnoDB
    DEFAULT CHARSET=utf8mb4
    COLLATE=utf8mb4_unicode_ci
    COMMENT='性別マスタ（コードのみ。表示名はi18n辞書で解決）';
