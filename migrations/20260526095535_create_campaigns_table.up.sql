-- Migration: create_campaigns_table
-- Created: 2026-05-26T09:55:35+04:00
-- Description: Tạo bảng campaigns và các index liên quan

BEGIN;

CREATE TABLE campaigns (
    id         UUID PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    code       VARCHAR(100) NOT NULL,
    status     VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at   TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE campaigns ADD CONSTRAINT uq_campaigns_code UNIQUE (code);

COMMIT;
