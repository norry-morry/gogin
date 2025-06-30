CREATE TABLE if not exists users (
--     `id`           CHAR(36)  not null comment 'ユーザーID',
    `id`           INT8 not null auto_increment comment 'ユーザーID',
    `name`         varchar(255)   not null comment 'ユーザー名-実名',
    `display_name` varchar(255)   not null comment '画面表示名',
    `email` varchar(255) DEFAULT NULL COMMENT 'email',
    `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT '作成日時',
    `updated_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT '更新日時',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '削除日時',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='アカウントマスタ';