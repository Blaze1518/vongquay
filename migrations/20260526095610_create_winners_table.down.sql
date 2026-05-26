-- Migration: create_winners_table (rollback)
-- Created: 2026-05-26T09:56:10+04:00

BEGIN;

DROP TABLE IF EXISTS winners;

COMMIT;
