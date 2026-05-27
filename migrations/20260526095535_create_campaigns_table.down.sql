-- Migration: create_campaigns_table (rollback)
-- Created: 2026-05-26T09:55:35+04:00

BEGIN;

DROP TABLE IF EXISTS campaigns;

COMMIT;
