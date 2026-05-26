-- Migration: create_winners_table
-- Created: 2026-05-26T09:56:10+04:00
-- Description: Tạo bảng winners lưu lịch sử trúng giải và các index liên quan

BEGIN;

CREATE TABLE winners (
    id          BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL,
    prize_id    BIGINT NOT NULL,
    ticket_id   BIGINT NOT NULL,
    draw_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tạo các index cho các trường ID quan hệ theo tag `gorm:"index"`
CREATE INDEX idx_winners_campaign_id ON winners(campaign_id);
CREATE INDEX idx_winners_prize_id ON winners(prize_id);
CREATE INDEX idx_winners_ticket_id ON winners(ticket_id);

COMMIT;
