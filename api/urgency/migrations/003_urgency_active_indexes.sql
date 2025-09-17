-- Migration: Add performance indexes for active urgencies
-- Date: 2025-09-17
-- Notes:
-- - These CREATE INDEX statements use CONCURRENTLY to avoid long locks.
-- - CONCURRENTLY cannot run inside a transaction; keep statements standalone.
-- - Safe to run multiple times thanks to IF NOT EXISTS.

-- 1) For listing active urgencies by sort priority asc then created_at desc
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_urgency_active_sort_created
    ON urgencies (sort_priority ASC, created_at DESC)
    WHERE deleted_at IS NULL;

-- 2) For listing assigned urgencies for an employee ordered by newest first
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_urgency_active_assignee_created
    ON urgencies (assigned_employee_id, created_at DESC)
    WHERE deleted_at IS NULL;

-- 3) For fetching only unassigned & active urgency IDs quickly (badge/counts)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_urgency_unassigned_active
    ON urgencies (id)
    WHERE deleted_at IS NULL AND sort_priority = 1;

