-- Migration: create_winners_table
-- Created: 2026-05-26T09:56:10+04:00
-- Description: Tạo bảng winners lưu lịch sử trúng giải và các index liên quan

BEGIN;

CREATE TABLE winners (
    id            UUID PRIMARY KEY,
    campaign_id   UUID NOT NULL,
    prize_id      UUID NOT NULL,
    prize_name    VARCHAR(255) NOT NULL,
    ticket_id     UUID NOT NULL,
    ticket_number VARCHAR(21) NOT NULL,
    username      VARCHAR(100) NOT NULL,
    draw_order    INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_winners_campaign_id FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE,
    CONSTRAINT fk_winners_prize_id FOREIGN KEY (prize_id) REFERENCES prizes(id) ON DELETE CASCADE,
    CONSTRAINT fk_winners_ticket_id FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,

    CONSTRAINT uq_winners_ticket_id UNIQUE (ticket_id)
);

CREATE INDEX idx_winners_campaign_id ON winners(campaign_id);
CREATE INDEX idx_winners_prize_id ON winners(prize_id);
CREATE INDEX idx_winners_username ON winners(username);

COMMIT;
