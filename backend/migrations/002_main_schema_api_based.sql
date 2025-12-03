-- EasyHire Main Database Schema - API Based
-- Version: 002
-- Created: $(date +"%Y-%m-%d %H:%M:%S")
-- Based on: docs/api/schemas/*.yaml

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================
-- 1. Users table (из Task #8)
-- ============================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'candidate', -- admin, hr, candidate, technical_expert
    company VARCHAR(255),
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- ============================
-- 2. Assessments table (строго по assessment.yaml)
-- ============================
CREATE TABLE assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id),
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    CHECK (status IN ('draft', 'published', 'active', 'completed', 'archived')),
    role VARCHAR(50),
    CHECK (role IN ('backend_developer', 'fullstack_developer', 'devops_engineer', 'team_lead')),
    target_level VARCHAR(50),
    CHECK (target_level IN ('junior', 'middle', 'senior', 'expert')),
    time_limit_minutes INTEGER CHECK (time_limit_minutes >= 1),
    question_count INTEGER CHECK (question_count >= 1),
    passing_score DECIMAL(5,2) CHECK (passing_score >= 0 AND passing_score <= 100),
    competency_weights JSONB DEFAULT '{}',
    level_distribution JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_assessments_created_by ON assessments(created_by);
CREATE INDEX idx_assessments_status ON assessments(status);
CREATE INDEX idx_assessments_role ON assessments(role);
CREATE INDEX idx_assessments_target_level ON assessments(target_level);

-- ============================
-- 3. Competencies table (упрощаем под API)
-- ============================
CREATE TABLE competencies (
    name VARCHAR(100) PRIMARY KEY, -- "concurrency", "system_design"
    category VARCHAR(50),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_competencies_category ON competencies(category);

-- ============================
-- 4. Assessment competencies (many-to-many)
-- ============================
CREATE TABLE assessment_competencies (
    assessment_id UUID REFERENCES assessments(id) ON DELETE CASCADE,
    competency_name VARCHAR(100) REFERENCES competencies(name) ON DELETE CASCADE,
    weight DECIMAL(5,2) DEFAULT 1.0,
    PRIMARY KEY (assessment_id, competency_name)
);

-- ============================
-- 5. Candidates table (строго по candidate.yaml)
-- ============================
CREATE TABLE candidates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(50),
    status VARCHAR(50) DEFAULT 'active',
    CHECK (status IN ('active', 'invited', 'completed', 'archived')),
    invited_at TIMESTAMP WITH TIME ZONE,
    completed_assessments_count INTEGER DEFAULT 0,
    average_score DECIMAL(5,2),
    last_assessment_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_candidates_email ON candidates(email);
CREATE INDEX idx_candidates_status ON candidates(status);
CREATE INDEX idx_candidates_invited_at ON candidates(invited_at);

-- ============================
-- 6. Questions table (строго по question.yaml)
-- ============================
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(50) NOT NULL,
    CHECK (type IN ('multiple_choice', 'coding', 'architecture', 'debugging')),
    competency VARCHAR(100) REFERENCES competencies(name),
    level VARCHAR(50) NOT NULL,
    CHECK (level IN ('junior', 'middle', 'senior', 'expert')),
    content JSONB NOT NULL, -- {text, code_snippet, options}
    test_cases JSONB DEFAULT '[]',
    solution JSONB, -- {code, explanation}
    explanation TEXT,
    ai_generated BOOLEAN DEFAULT FALSE,
    validation_status VARCHAR(50) DEFAULT 'pending',
    CHECK (validation_status IN ('pending', 'approved', 'rejected', 'needs_review')),
    validated_by UUID REFERENCES users(id),
    validated_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_questions_competency ON questions(competency);
CREATE INDEX idx_questions_level ON questions(level);
CREATE INDEX idx_questions_type ON questions(type);
CREATE INDEX idx_questions_validation_status ON questions(validation_status);

-- ============================
-- 7. Assessment questions (many-to-many)
-- ============================
CREATE TABLE assessment_questions (
    assessment_id UUID REFERENCES assessments(id) ON DELETE CASCADE,
    question_id UUID REFERENCES questions(id) ON DELETE CASCADE,
    order_index INTEGER DEFAULT 0,
    points INTEGER DEFAULT 10,
    PRIMARY KEY (assessment_id, question_id)
);

CREATE INDEX idx_assessment_questions_order ON assessment_questions(order_index);

-- ============================
-- 8. Results table (строго по result.yaml)
-- ============================
CREATE TABLE results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assessment_id UUID REFERENCES assessments(id),
    candidate_id UUID REFERENCES candidates(id),
    score DECIMAL(5,2) NOT NULL,
    CHECK (score >= 0 AND score <= 100),
    fibonacci_score DECIMAL(5,2),
    level VARCHAR(50) NOT NULL,
    CHECK (level IN ('trainee', 'junior', 'middle', 'senior', 'expert')),
    passed BOOLEAN NOT NULL,
    time_spent_minutes INTEGER,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    competency_breakdown JSONB DEFAULT '[]',
    question_results JSONB DEFAULT '[]',
    feedback JSONB DEFAULT '[]',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_results_assessment_id ON results(assessment_id);
CREATE INDEX idx_results_candidate_id ON results(candidate_id);
CREATE INDEX idx_results_score ON results(score);
CREATE INDEX idx_results_level ON results(level);
CREATE INDEX idx_results_passed ON results(passed);
CREATE INDEX idx_results_created_at ON results(created_at);

-- ============================
-- 9. Invitations table (для приглашений кандидатов)
-- ============================
CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assessment_id UUID REFERENCES assessments(id),
    candidate_id UUID REFERENCES candidates(id),
    invited_by UUID REFERENCES users(id),
    token VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(50) DEFAULT 'sent',
    CHECK (status IN ('sent', 'opened', 'started', 'completed', 'expired')),
    email_sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_invitations_token ON invitations(token);
CREATE INDEX idx_invitations_assessment_id ON invitations(assessment_id);
CREATE INDEX idx_invitations_candidate_id ON invitations(candidate_id);
CREATE INDEX idx_invitations_status ON invitations(status);

-- ============================
-- 10. Code executions table (для выполнения кода)
-- ============================
CREATE TABLE code_executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    result_id UUID REFERENCES results(id) ON DELETE CASCADE,
    question_id UUID REFERENCES questions(id),
    candidate_code TEXT NOT NULL,
    language VARCHAR(50) DEFAULT 'go',
    status VARCHAR(50) NOT NULL,
    CHECK (status IN ('pending', 'running', 'success', 'error', 'timeout')),
    exit_code INTEGER,
    stdout TEXT,
    stderr TEXT,
    execution_time_ms INTEGER,
    memory_used_kb INTEGER,
    docker_container_id VARCHAR(100),
    error_message TEXT,
    logs TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_code_executions_result_id ON code_executions(result_id);
CREATE INDEX idx_code_executions_status ON code_executions(status);

-- ============================
-- Update schema migrations
-- ============================
INSERT INTO schema_migrations (version, name) 
VALUES (2, 'main_schema_api_based')
ON CONFLICT (version) DO NOTHING;

-- ============================
-- Seed initial competencies (из API схем)
-- ============================
INSERT INTO competencies (name, category, description) VALUES
('concurrency', 'backend', 'Goroutines, channels, synchronization'),
('system_design', 'backend', 'Architecture patterns, scalability'),
('testing', 'backend', 'Unit tests, integration tests'),
('microservices', 'backend', 'Distributed systems, service communication'),
('database_design', 'backend', 'SQL, ORM, optimization'),
('algorithms', 'computer-science', 'Data structures and algorithms'),
('security', 'backend', 'Authentication, authorization, encryption'),
('debugging', 'backend', 'Troubleshooting, profiling'),
('architecture', 'backend', 'System design, patterns')
ON CONFLICT (name) DO NOTHING;

-- Create default admin user
INSERT INTO users (email, password_hash, name, role) 
VALUES (
    'admin@easyhire.dev', 
    crypt('admin123', gen_salt('bf')), 
    'System Administrator', 
    'admin'
) ON CONFLICT (email) DO NOTHING;

-- Create default hr user
INSERT INTO users (email, password_hash, name, role) 
VALUES (
    'hr@company.com', 
    crypt('hr123456', gen_salt('bf')), 
    'HR Manager', 
    'hr'
) ON CONFLICT (email) DO NOTHING;
