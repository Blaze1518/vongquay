-- Migration: update_tickets_table_add_username (rollback)
-- Created: 2026-05-26T17:53:45+04:00

BEGIN;

DROP INDEX IF EXISTS idx_tickets_campaign_username;
DROP INDEX IF EXISTS idx_tickets_campaign_code;

CREATE UNIQUE INDEX idx_tickets_code ON tickets(ticket_number);

ALTER TABLE tickets DROP COLUMN IF EXISTS is_canceled;
ALTER TABLE tickets DROP COLUMN IF EXISTS username;

COMMIT;
