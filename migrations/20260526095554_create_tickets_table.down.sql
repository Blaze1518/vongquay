-- Migration: create_tickets_table (rollback)
-- Created: 2026-05-26T09:55:54+04:00

BEGIN;

DROP TABLE IF EXISTS tickets;

COMMIT;
