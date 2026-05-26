-- Migration: create_tickets_table
-- Created: 2026-05-26T09:55:54+04:00
-- Description: Tạo bảng tickets và index liên quan

BEGIN;

CREATE TABLE tickets (
    id            BIGSERIAL PRIMARY KEY,
    campaign_id   BIGINT NOT NULL,
    ticket_number VARCHAR(21) NOT NULL,
    username      VARCHAR(100) NOT NULL,
    is_winner     BOOLEAN NOT NULL DEFAULT FALSE,
    is_canceled   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tickets_campaign_username ON tickets(campaign_id, username);
CREATE UNIQUE INDEX idx_tickets_code ON tickets(ticket_number);

COMMIT;
