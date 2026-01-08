CREATE TABLE user_addresses (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT 'ユーザー住所ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT 'ユーザーID（users.id）',
    purpose_id BIGINT UNSIGNED NOT NULL COMMENT '住所用途ID（address_purposes.id）',
    is_primary BOOLEAN NOT NULL DEFAULT FALSE COMMENT '代表住所フラグ（同一ユーザー×用途で最大1件）',

    country_code CHAR(2) NOT NULL COMMENT '国コード（ISO 3166-1 alpha-2: JP/US 等）',
    administrative_area VARCHAR(128) NULL COMMENT '第一行政区画（都道府県/州など）',
    locality           VARCHAR(128) NULL COMMENT '第二行政区画（市区町村など）',
    dependent_locality VARCHAR(128) NULL COMMENT '第三行政区画（区/町域/地区など任意）',
    postal_code        VARCHAR(32)  NULL COMMENT '郵便番号（国により形式可変）',
    sorting_code       VARCHAR(32)  NULL COMMENT '仕分けコード（一部の国で使用）',

    address_line1      VARCHAR(160) NOT NULL COMMENT '住所行1（番地/丁目/通り等）',
    address_line2      VARCHAR(160) NULL COMMENT '住所行2（建物名/部屋番号等）',
    address_line3      VARCHAR(160) NULL COMMENT '住所行3（予備）',

    organization       VARCHAR(160) NULL COMMENT '宛先：会社/団体名（任意）',
    given_name         VARCHAR(80)  NULL COMMENT '宛名：名（任意）',
    family_name        VARCHAR(80)  NULL COMMENT '宛名：姓（任意）',

    language_code      VARCHAR(35)  NULL COMMENT '言語タグ（BCP47: ja-JP/en-US 等）',
    latitude           DECIMAL(9,6) NULL COMMENT '緯度',
    longitude          DECIMAL(9,6) NULL COMMENT '経度',

    created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    deleted_at         TIMESTAMP NULL COMMENT '論理削除',

    -- 有効な代表住所のみ一意にしたいので、deleted_at IS NULL も条件に含める
    active_primary_flag TINYINT
        GENERATED ALWAYS AS (IF(is_primary AND deleted_at IS NULL, 1, NULL))
        STORED COMMENT '有効な代表一意制約用の生成列',

    CONSTRAINT fk_user_addresses_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_addresses_purpose
        FOREIGN KEY (purpose_id) REFERENCES address_purposes(id) ON DELETE RESTRICT
)
    ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci
COMMENT='ユーザー住所（国際対応の抽象スキーマ：行政区＋可変住所行）';

-- 代表住所：ユーザー×用途で最大1件（is_primary=true の行が1つだけ許容）
CREATE UNIQUE INDEX uq_user_purpose_primary
    ON user_addresses (user_id, purpose_id, active_primary_flag);

-- よく使う検索系
CREATE INDEX idx_user_addresses_user
    ON user_addresses (user_id, is_primary);

CREATE INDEX idx_user_addresses_country_postal
    ON user_addresses (country_code, postal_code);

CREATE INDEX idx_user_addresses_loc
    ON user_addresses (country_code, administrative_area, locality);
