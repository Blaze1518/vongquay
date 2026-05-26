-- Migration: create_campaigns_table (rollback)
-- Created: 2026-05-26T09:55:35+04:00

BEGIN;

-- Xóa bảng campaigns (Index idx_campaigns_code sẽ tự động bị xóa theo bảng)
DROP TABLE IF EXISTS campaigns;

COMMIT;
