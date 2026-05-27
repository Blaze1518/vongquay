-- Migration: create_ticketimportjobs_table
-- Created: 2026-05-26T17:40:42+04:00
-- Description: Tạo bảng ticketimportjobs phục vụ lưu tiến độ import file 100k dòng

BEGIN;

CREATE TABLE ticket_import_jobs (
    id           UUID PRIMARY KEY,
    campaign_id  UUID NOT NULL,
    status       VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    total_rows   INT NOT NULL DEFAULT 0,
    success_rows INT NOT NULL DEFAULT 0,
    failed_rows  INT NOT NULL DEFAULT 0,
    error_log    TEXT,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_ticket_import_jobs_campaign_id FOREIGN KEY (campaign_id) 
        REFERENCES campaigns(id) ON DELETE CASCADE
);

CREATE INDEX idx_ticket_import_jobs_campaign_id ON ticket_import_jobs(campaign_id);

CREATE INDEX idx_ticket_import_jobs_status ON ticket_import_jobs(status) 
WHERE status IN ('PENDING', 'PROCESSING');

COMMIT;
