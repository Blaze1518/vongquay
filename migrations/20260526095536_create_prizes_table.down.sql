-- Migration: create_prizes_table (rollback)
-- Created: 2026-05-26T09:54:22+04:00

BEGIN;

DROP TABLE IF EXISTS prizes;

COMMIT;
