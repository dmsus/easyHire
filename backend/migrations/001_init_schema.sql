-- Initial schema for EasyHire
-- Version: 001
-- Created: $(date)

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Schema migration tracking
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO schema_migrations (version, name) 
VALUES (1, 'initial_schema')
ON CONFLICT (version) DO NOTHING;
