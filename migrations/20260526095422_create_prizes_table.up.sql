-- Migration: create_prizes_table
-- Created: 2026-05-26T09:54:22+04:00
-- Description: Tạo bảng prizes và index liên quan

BEGIN;

CREATE TABLE prizes (
    id          BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL,
    name        VARCHAR(255) NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    priority    INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_prizes_campaign_id ON prizes(campaign_id);

COMMIT;
