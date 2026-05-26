-- Migration: create_ticketimportjobs_table (rollback)
-- Created: 2026-05-26T17:40:42+04:00

BEGIN;

DROP TABLE IF EXISTS ticketimportjobs;

COMMIT;
