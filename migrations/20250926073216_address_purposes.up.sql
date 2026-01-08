CREATE TABLE address_purposes (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '住所用途ID（マスター）',
    code VARCHAR(64) NOT NULL UNIQUE COMMENT '用途コード（例: shipping/billing/home/office/other）',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '表示順',
    is_active TINYINT(1) NOT NULL DEFAULT 1 COMMENT '有効フラグ（1:true/0:false）',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時'
)
    ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci
COMMENT='住所用途マスタ（用途の種類を拡張可能に管理）';
