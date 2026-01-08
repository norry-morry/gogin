CREATE TABLE `user_identities` (
   `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
   `user_id`           BIGINT UNSIGNED NOT NULL COMMENT 'users.id',
   `provider`          VARCHAR(64)     NOT NULL COMMENT 'google.com, github.com, password, ...',
   `provider_user_id`  VARCHAR(255)    NOT NULL COMMENT 'IdP 側恒久ID（sub/subject など）',
   `provider_display_name`  VARCHAR(255)    NOT NULL COMMENT '連携時氏名のスナップショット（監査用）',
   `email_at_signup`   VARCHAR(320)    NULL COMMENT '連携時メールのスナップショット（監査用）',
   `created_at`        DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

   PRIMARY KEY (`id`),
   UNIQUE KEY `uq_user_identity_provider_user` (`provider`, `provider_user_id`),
   KEY `idx_user_identities_user` (`user_id`),

   CONSTRAINT `fk_auth_identities_user`
       FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
           ON DELETE CASCADE
           ON UPDATE RESTRICT
)
    ENGINE=InnoDB
    DEFAULT CHARSET=utf8mb4
    COLLATE=utf8mb4_unicode_ci
    COMMENT='プロバイダデータ';
