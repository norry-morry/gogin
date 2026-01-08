CREATE TABLE `users` (
     `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
     `uid`             VARCHAR(128)    NOT NULL COMMENT 'Firebase UID（プロジェクト内一意）',
     `email`           VARCHAR(320)    NULL COMMENT 'email',
     `email_verified`  TINYINT(1)      NOT NULL DEFAULT 0 COMMENT 'メールアドレスの承認',
     `display_name`    VARCHAR(255)    NULL COMMENT '画面表示名(正)',
     `photo_url`       TEXT            NULL COMMENT 'アバターアイコン',
     `disabled`        TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '一時停止フラグ',
     `deleted_at`      DATETIME(6)     NULL COMMENT '論理削除（退会）',
     `last_login_at`   DATETIME(6)     NULL COMMENT '最終ログイン日時',
     `created_at`      DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
     `updated_at`      DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
     PRIMARY KEY (`id`),
     UNIQUE KEY `uq_users_uid`   (`uid`),
     UNIQUE KEY `uq_users_email` (`email`),          -- MySQL の UNIQUE は NULL を複数許容
     KEY `idx_users_deleted_at`  (`deleted_at`)      -- 退会除外のフィルタ用
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='アカウントマスタ';
