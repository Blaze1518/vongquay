-- Migration: update_tickets_table_add_username
-- Created: 2026-05-26T17:53:45+04:00
-- Description: Add description here

BEGIN;

ALTER TABLE tickets ADD COLUMN username VARCHAR(100) NOT NULL;
ALTER TABLE tickets ADD COLUMN is_canceled BOOLEAN NOT NULL DEFAULT FALSE;

DROP INDEX IF EXISTS idx_tickets_code;

CREATE UNIQUE INDEX idx_tickets_campaign_code ON tickets(campaign_id, ticket_number);

CREATE INDEX idx_tickets_campaign_username ON tickets(campaign_id, username);

CREATE INDEX idx_tickets_available_for_draw 
ON tickets (campaign_id) 
WHERE is_winner = FALSE AND is_canceled = FALSE;

COMMIT;
