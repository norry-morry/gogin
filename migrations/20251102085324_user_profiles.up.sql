CREATE TABLE user_profiles (
    user_id BIGINT UNSIGNED NOT NULL COMMENT 'users.id と1:1',
    -- 氏名（必須）
    family_name VARCHAR(80) NOT NULL COMMENT '姓',
    given_name  VARCHAR(80)  NOT NULL COMMENT '名',

    -- カナ（必須）
    family_name_kana VARCHAR(80) NOT NULL COMMENT 'セイ',
    given_name_kana  VARCHAR(80) NOT NULL COMMENT 'メイ',

    -- 表示用の結合（常に姓＋空白＋名）。アプリ側で扱う “印字用” 正規値として利用
    legal_name VARCHAR(161)
        AS (CONCAT_WS(' ', family_name, given_name))
    STORED
    COMMENT '表示・帳票用の結合済み氏名',

    -- カナの結合（常にセイ＋空白＋メイ）。アプリ側で扱う “印字用” 正規値として利用
    legal_name_kana VARCHAR(161)
        AS (CONCAT_WS(' ', family_name_kana, given_name_kana))
    STORED
    COMMENT '結合済みカナ',

    birth_date DATE NULL COMMENT '生年月日',

    -- ★ 名称を gender_id にし、FK先が明確に
    gender_id TINYINT UNSIGNED NULL COMMENT 'genders.id',

    initial varchar(32) NULL COMMENT 'イニシャル',

    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',

    PRIMARY KEY (user_id),

    -- ★ users への 1:1 を厳格化（CASCADE はOK）
    CONSTRAINT fk_user_profiles_user
        FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE
            ON UPDATE RESTRICT,

    -- ★ genders への FK は RESTRICT（マスタ削除で子を消さない）
    CONSTRAINT fk_user_profiles_gender
        FOREIGN KEY (gender_id) REFERENCES genders(id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
)
    ENGINE=InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci
    COMMENT='ユーザー基本情報（1ユーザー1行）';

    -- よく使いそうな検索（名字・印字名・カナ）
    CREATE INDEX idx_user_profiles_family ON user_profiles (family_name);
    CREATE INDEX idx_user_profiles_legal  ON user_profiles (legal_name);
