-- Migration: Add metadata column to users table
-- Description: Adds JSONB metadata column for storing additional user information

BEGIN;

-- Add metadata column to users table
ALTER TABLE users 
ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;

-- Add comment for documentation
COMMENT ON COLUMN users.metadata IS 'Additional user metadata stored as JSON';

COMMIT;
