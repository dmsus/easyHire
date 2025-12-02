-- Scoring system enhancements
-- Version: 004
-- Created: $(date +"%Y-%m-%d %H:%M:%S")

-- Add base_weight to competencies
ALTER TABLE competencies ADD COLUMN IF NOT EXISTS base_weight DECIMAL(3,2) DEFAULT 1.0;

-- Update existing competencies with weights from the matrix
UPDATE competencies SET base_weight = 1.3 WHERE name = 'concurrency';
UPDATE competencies SET base_weight = 1.0 WHERE name = 'go_fundamentals';
UPDATE competencies SET base_weight = 1.1 WHERE name = 'data_structures_go';
UPDATE competencies SET base_weight = 1.1 WHERE name = 'memory_management';
UPDATE competencies SET base_weight = 1.0 WHERE name = 'http_go';
UPDATE competencies SET base_weight = 1.3 WHERE name = 'system_design';
UPDATE competencies SET base_weight = 1.2 WHERE name = 'microservices';
UPDATE competencies SET base_weight = 1.1 WHERE name = 'containerization';
UPDATE competencies SET base_weight = 1.2 WHERE name = 'reliability';
UPDATE competencies SET base_weight = 1.2 WHERE name = 'message_brokers';
UPDATE competencies SET base_weight = 1.2 WHERE name = 'software_design';
UPDATE competencies SET base_weight = 1.3 WHERE name = 'architecture';
UPDATE competencies SET base_weight = 1.1 WHERE name = 'quality_assurance';
UPDATE competencies SET base_weight = 1.1 WHERE name = 'optimization';
UPDATE competencies SET base_weight = 1.0 WHERE name = 'sdlc';
UPDATE competencies SET base_weight = 1.0 WHERE name = 'requirements';
UPDATE competencies SET base_weight = 1.1 WHERE name = 'ci_cd';
UPDATE competencies SET base_weight = 1.0 WHERE name = 'git';
UPDATE competencies SET base_weight = 1.2 WHERE name = 'web_security';
UPDATE competencies SET base_weight = 1.3 WHERE name = 'data_security';

-- Create table for level weights (Fibonacci)
CREATE TABLE IF NOT EXISTS level_weights (
    level VARCHAR(20) PRIMARY KEY,
    weight INTEGER NOT NULL
);

INSERT INTO level_weights (level, weight) VALUES
('junior', 1),
('middle', 2),
('senior', 3),
('expert', 5)
ON CONFLICT (level) DO NOTHING;

-- Create table for level thresholds (Fibonacci)
CREATE TABLE IF NOT EXISTS level_thresholds (
    level VARCHAR(20) PRIMARY KEY,
    threshold INTEGER NOT NULL
);

INSERT INTO level_thresholds (level, threshold) VALUES
('junior', 8),
('middle', 21),
('senior', 55),
('expert', 144)
ON CONFLICT (level) DO NOTHING;

-- Update schema migrations
INSERT INTO schema_migrations (version, name) 
VALUES (4, 'scoring_system')
ON CONFLICT (version) DO NOTHING;
