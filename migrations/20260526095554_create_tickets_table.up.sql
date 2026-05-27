-- Migration: create_tickets_table
-- Created: 2026-05-26T09:55:54+04:00
-- Description: Tạo bảng tickets và index liên quan

BEGIN;

CREATE TABLE tickets (
    id            UUID PRIMARY KEY,
    campaign_id   UUID NOT NULL,
    ticket_number VARCHAR(21) NOT NULL,
    username      VARCHAR(100) NOT NULL,
    is_winner     BOOLEAN NOT NULL DEFAULT FALSE,
    is_canceled   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_tickets_campaign_id FOREIGN KEY (campaign_id) 
        REFERENCES campaigns(id) ON DELETE CASCADE
);

-- Ràng buộc: Trong một chiến dịch, một số vé (ticket_number) chỉ được xuất hiện một lần
ALTER TABLE tickets ADD CONSTRAINT uq_tickets_campaign_ticket UNIQUE (campaign_id, ticket_number);

-- Index tăng tốc khi tìm kiếm vé theo người dùng trong một chiến dịch cụ thể
CREATE INDEX idx_tickets_campaign_username ON tickets(campaign_id, username);

CREATE INDEX idx_tickets_available_for_draw 
ON tickets (campaign_id) 
WHERE is_winner = FALSE AND is_canceled = FALSE;

COMMIT;
