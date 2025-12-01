-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pgvector";

-- Create schema if needed
CREATE SCHEMA IF NOT EXISTS easyhire;

-- Set search path
SET search_path TO easyhire, public;
