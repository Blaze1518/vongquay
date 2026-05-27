-- Migration: add_fields_to_winners_table
-- Created: 2026-05-27T09:48:02+04:00
-- Description: Add description here

BEGIN;

ALTER TABLE winners ADD COLUMN ticket_number VARCHAR(21) NOT NULL;
ALTER TABLE winners ADD COLUMN username VARCHAR(100) NOT NULL;

CREATE INDEX idx_winners_username ON winners(username);

COMMIT;
