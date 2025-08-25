-- Migration: Simplify Activity Model
-- Description: Remove complex fields and simplify activity table structure
-- Date: 2025-08-25

-- End any existing transaction first
ROLLBACK;

-- Start fresh transaction
BEGIN;

-- Skip backup creation if already exists (uncomment next lines if you want a fresh backup)
-- DROP TABLE IF EXISTS activities_backup;
-- CREATE TABLE activities_backup AS SELECT * FROM activities;

-- Add new columns for simplified model (only if they don't exist)
ALTER TABLE activities
ADD COLUMN IF NOT EXISTS employee_id INTEGER,
ADD COLUMN IF NOT EXISTS urgency_id INTEGER;

-- Migrate existing data (only if old columns exist)
DO $$
BEGIN
    -- Check if old columns exist before trying to migrate data
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'activities' AND column_name = 'actor_id')
       AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'activities' AND column_name = 'target_id') THEN

        -- Migrate data from old columns to new columns
        UPDATE activities
        SET employee_id = COALESCE(actor_id, 0),
            urgency_id = COALESCE(target_id, 0)
        WHERE (EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'activities' AND column_name = 'target_type')
               AND (target_type = 'urgency' OR target_type IS NULL))
           OR NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'activities' AND column_name = 'target_type');

        -- Set default values for any remaining NULL values
        UPDATE activities
        SET employee_id = COALESCE(employee_id, 1),
            urgency_id = COALESCE(urgency_id, 1)
        WHERE employee_id IS NULL OR urgency_id IS NULL;

    END IF;
END $$;

-- Make the new columns NOT NULL after data migration (only if they exist and have data)
DO $$
BEGIN
    -- Check if columns exist and set NOT NULL
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'activities' AND column_name = 'employee_id') THEN
        ALTER TABLE activities ALTER COLUMN employee_id SET NOT NULL;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'activities' AND column_name = 'urgency_id') THEN
        ALTER TABLE activities ALTER COLUMN urgency_id SET NOT NULL;
    END IF;
END $$;

-- Remove old complex columns including level (activities don't have levels, only urgencies do)
ALTER TABLE activities
DROP COLUMN IF EXISTS type,
DROP COLUMN IF EXISTS level,
DROP COLUMN IF EXISTS title,
DROP COLUMN IF EXISTS actor_id,
DROP COLUMN IF EXISTS actor_name,
DROP COLUMN IF EXISTS target_id,
DROP COLUMN IF EXISTS target_type,
DROP COLUMN IF EXISTS metadata;

-- Add indexes for better performance (only if columns exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'activities' AND column_name = 'employee_id') THEN
        CREATE INDEX IF NOT EXISTS idx_activities_employee_id ON activities(employee_id);
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'activities' AND column_name = 'urgency_id') THEN
        CREATE INDEX IF NOT EXISTS idx_activities_urgency_id ON activities(urgency_id);
    END IF;

    CREATE INDEX IF NOT EXISTS idx_activities_created_at ON activities(created_at);
END $$;

-- Add foreign key constraints (optional, depends on your setup)
-- Uncomment these if you want referential integrity
-- ALTER TABLE activities 
-- ADD CONSTRAINT fk_activities_employee_id 
-- FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE;

-- ALTER TABLE activities 
-- ADD CONSTRAINT fk_activities_urgency_id 
-- FOREIGN KEY (urgency_id) REFERENCES urgencies(id) ON DELETE CASCADE;

-- No level constraints needed since activities don't have levels

-- Commit transaction
COMMIT;

-- Drop backup table after successful migration (optional)
-- DROP TABLE activities_backup;
