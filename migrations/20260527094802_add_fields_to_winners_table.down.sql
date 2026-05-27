-- Migration: add_fields_to_winners_table (rollback)
-- Created: 2026-05-27T09:48:02+04:00

BEGIN;

DROP INDEX IF EXISTS idx_winners_username;

ALTER TABLE winners DROP COLUMN IF EXISTS username;
ALTER TABLE winners DROP COLUMN IF EXISTS ticket_number;

COMMIT;
