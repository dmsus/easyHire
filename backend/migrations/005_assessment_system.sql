-- Assessment System Migration
BEGIN;

-- Create enum types for assessment system
DO $$ BEGIN
    CREATE TYPE assessment_status AS ENUM ('draft', 'active', 'paused', 'archived');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE assessment_type AS ENUM ('screening', 'technical', 'coding', 'system_design', 'custom');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE session_status AS ENUM ('pending', 'in_progress', 'completed', 'expired', 'terminated');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE question_type AS ENUM ('multiple_choice', 'coding', 'architecture', 'debugging');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE question_status AS ENUM ('draft', 'active', 'inactive', 'archived');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE invitation_status AS ENUM ('pending', 'sent', 'opened', 'accepted', 'completed', 'expired');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE difficulty_level AS ENUM ('junior', 'middle', 'senior', 'expert');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Assessments table
CREATE TABLE IF NOT EXISTS assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status assessment_status NOT NULL DEFAULT 'draft',
    type assessment_type NOT NULL,
    target_level difficulty_level NOT NULL,
    time_limit INTEGER NOT NULL DEFAULT 3600,
    total_questions INTEGER NOT NULL DEFAULT 20,
    passing_score DECIMAL(5,2) DEFAULT 70.0,
    shuffle_questions BOOLEAN DEFAULT TRUE,
    show_explanation BOOLEAN DEFAULT FALSE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Assessment competencies (many-to-many with weights)
CREATE TABLE IF NOT EXISTS assessment_competencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    weight DECIMAL(3,2) NOT NULL DEFAULT 1.0,
    min_questions INTEGER NOT NULL DEFAULT 1,
    max_questions INTEGER NOT NULL DEFAULT 5,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assessment_id, competency_id)
);

-- Questions table
CREATE TABLE IF NOT EXISTS questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type question_type NOT NULL,
    competency_id UUID NOT NULL REFERENCES competencies(id) ON DELETE CASCADE,
    level difficulty_level NOT NULL,
    code_template TEXT,
    test_cases TEXT,
    options JSONB,
    correct_answer TEXT,
    explanation TEXT,
    estimated_time INTEGER NOT NULL DEFAULT 180,
    points DECIMAL(5,2) NOT NULL DEFAULT 1.0,
    status question_status NOT NULL DEFAULT 'active',
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Question tags
CREATE TABLE IF NOT EXISTS question_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Question-tag many-to-many
CREATE TABLE IF NOT EXISTS question_tag_associations (
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES question_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (question_id, tag_id)
);

-- Assessment questions (with ordering and versioning)
CREATE TABLE IF NOT EXISTS assessment_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    "order" INTEGER NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    is_required BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assessment_id, question_id)
);

-- Assessment sessions
CREATE TABLE IF NOT EXISTS assessment_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    candidate_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status session_status NOT NULL DEFAULT 'pending',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    time_spent INTEGER DEFAULT 0,
    score DECIMAL(5,2),
    percentage DECIMAL(5,2),
    level difficulty_level,
    cheating_attempts INTEGER DEFAULT 0,
    browser_info JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assessment_id, candidate_id) WHERE status NOT IN ('completed', 'expired', 'terminated')
);

-- Candidate answers
CREATE TABLE IF NOT EXISTS candidate_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES assessment_sessions(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    answer TEXT,
    code TEXT,
    is_correct BOOLEAN,
    points_earned DECIMAL(5,2),
    time_spent INTEGER,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMP,
    cheating_flags JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id, question_id)
);

-- Invitations
CREATE TABLE IF NOT EXISTS invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    candidate_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(64) UNIQUE NOT NULL,
    status invitation_status NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    opened_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Assessment tags
CREATE TABLE IF NOT EXISTS assessment_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Assessment-tag many-to-many
CREATE TABLE IF NOT EXISTS assessment_tag_associations (
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES assessment_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (assessment_id, tag_id)
);

-- Results table
CREATE TABLE IF NOT EXISTS results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES assessment_sessions(id) ON DELETE CASCADE UNIQUE,
    total_score DECIMAL(5,2) NOT NULL,
    percentage DECIMAL(5,2) NOT NULL,
    level difficulty_level NOT NULL,
    time_spent INTEGER NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    competency_breakdown JSONB,
    recommendations TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_assessments_created_by ON assessments(created_by);
CREATE INDEX IF NOT EXISTS idx_assessments_status ON assessments(status);
CREATE INDEX IF NOT EXISTS idx_assessments_target_level ON assessments(target_level);
CREATE INDEX IF NOT EXISTS idx_assessment_competencies_assessment ON assessment_competencies(assessment_id);
CREATE INDEX IF NOT EXISTS idx_assessment_competencies_competency ON assessment_competencies(competency_id);
CREATE INDEX IF NOT EXISTS idx_questions_competency ON questions(competency_id);
CREATE INDEX IF NOT EXISTS idx_questions_level ON questions(level);
CREATE INDEX IF NOT EXISTS idx_questions_status ON questions(status);
CREATE INDEX IF NOT EXISTS idx_assessment_questions_assessment ON assessment_questions(assessment_id);
CREATE INDEX IF NOT EXISTS idx_assessment_questions_question ON assessment_questions(question_id);
CREATE INDEX IF NOT EXISTS idx_assessment_sessions_candidate ON assessment_sessions(candidate_id);
CREATE INDEX IF NOT EXISTS idx_assessment_sessions_status ON assessment_sessions(status);
CREATE INDEX IF NOT EXISTS idx_assessment_sessions_assessment ON assessment_sessions(assessment_id);
CREATE INDEX IF NOT EXISTS idx_candidate_answers_session ON candidate_answers(session_id);
CREATE INDEX IF NOT EXISTS idx_candidate_answers_question ON candidate_answers(question_id);
CREATE INDEX IF NOT EXISTS idx_invitations_token ON invitations(token);
CREATE INDEX IF NOT EXISTS idx_invitations_assessment ON invitations(assessment_id);
CREATE INDEX IF NOT EXISTS idx_invitations_candidate ON invitations(candidate_id);
CREATE INDEX IF NOT EXISTS idx_results_session ON results(session_id);

COMMIT;
