-- Seed data for assessment system
BEGIN;

-- Insert sample competencies if they don't exist
INSERT INTO competencies (id, name, description, category, weight, created_at)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'go_fundamentals', 'Go Fundamentals', 'core_go', 1.0, NOW()),
    ('22222222-2222-2222-2222-222222222222', 'concurrency', 'Concurrency', 'core_go', 1.3, NOW()),
    ('33333333-3333-3333-3333-333333333333', 'data_structures_go', 'Data Structures in Go', 'core_go', 1.1, NOW()),
    ('44444444-4444-4444-4444-444444444444', 'system_design', 'System Design', 'architecture', 1.3, NOW()),
    ('55555555-5555-5555-5555-555555555555', 'microservices', 'Microservices', 'architecture', 1.2, NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample question tags
INSERT INTO question_tags (id, name, description, color, created_at)
VALUES 
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'algorithm', 'Algorithm questions', '#FF6B6B', NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'database', 'Database related questions', '#4ECDC4', NOW()),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'security', 'Security questions', '#FFD166', NOW()),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'performance', 'Performance optimization', '#06D6A0', NOW()),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'testing', 'Testing and QA', '#118AB2', NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample questions
INSERT INTO questions (id, title, description, type, competency_id, level, code_template, test_cases, options, correct_answer, explanation, estimated_time, points, status, created_by, created_at)
VALUES 
    -- Go Fundamentals questions
    ('aaaaaaaa-0001-aaaa-aaaa-aaaaaaaaaaaa', 'Go Variables and Types', 'What is the zero value of a boolean variable in Go?', 'multiple_choice', '11111111-1111-1111-1111-111111111111', 'junior', '', '', '["true", "false", "nil", "0"]', 'false', 'In Go, the zero value for boolean type is false.', 60, 1.0, 'active', NULL, NOW()),
    
    ('aaaaaaaa-0002-aaaa-aaaa-aaaaaaaaaaaa', 'Go Functions', 'Write a function that returns the sum of two integers', 'coding', '11111111-1111-1111-1111-111111111111', 'junior', 'func sum(a int, b int) int {\n    // Your code here\n}', '[{"input": "1, 2", "output": "3"}, {"input": "0, 0", "output": "0"}, {"input": "-1, 1", "output": "0"}]', NULL, 'func sum(a int, b int) int {\n    return a + b\n}', 'Basic function implementation in Go.', 120, 1.0, 'active', NULL, NOW()),
    
    -- Concurrency questions
    ('bbbbbbbb-0001-bbbb-bbbb-bbbbbbbbbbbb', 'Goroutines', 'What happens when you start a goroutine?', 'multiple_choice', '22222222-2222-2222-2222-222222222222', 'middle', '', '', '["It blocks the main thread", "It runs concurrently", "It must complete before main exits", "It requires explicit thread management"]', 'It runs concurrently', 'Goroutines are lightweight threads managed by the Go runtime.', 90, 2.0, 'active', NULL, NOW()),
    
    ('bbbbbbbb-0002-bbbb-bbbb-bbbbbbbbbbbb', 'WaitGroup Usage', 'Implement a function that uses sync.WaitGroup to wait for multiple goroutines', 'coding', '22222222-2222-2222-2222-222222222222', 'middle', 'package main\n\nimport (\n    "fmt"\n    "sync"\n)\n\nfunc process(id int, wg *sync.WaitGroup) {\n    defer wg.Done()\n    // Your code here\n}\n\nfunc main() {\n    var wg sync.WaitGroup\n    // Your code here\n}', '[{"input": "", "output": "Processing 1\\nProcessing 2\\nProcessing 3\\nAll goroutines completed"}]', NULL, 'package main\n\nimport (\n    "fmt"\n    "sync"\n)\n\nfunc process(id int, wg *sync.WaitGroup) {\n    defer wg.Done()\n    fmt.Printf("Processing %d\\n", id)\n}\n\nfunc main() {\n    var wg sync.WaitGroup\n    \n    for i := 1; i <= 3; i++ {\n        wg.Add(1)\n        go process(i, &wg)\n    }\n    \n    wg.Wait()\n    fmt.Println("All goroutines completed")\n}', 'Proper use of WaitGroup for goroutine synchronization.', 180, 2.0, 'active', NULL, NOW()),
    
    -- System Design questions
    ('cccccccc-0001-cccc-cccc-cccccccccccc', 'Load Balancer Design', 'Design a load balancer for a web service', 'architecture', '44444444-4444-4444-4444-444444444444', 'senior', '', '', NULL, 'Design should include: 1) Multiple server instances, 2) Load balancing algorithm (round-robin, least connections), 3) Health checks, 4) Session persistence if needed, 5) Scalability considerations.', 'Designing load balancers for high availability.', 300, 3.0, 'active', NULL, NOW())
ON CONFLICT (id) DO NOTHING;

-- Associate questions with tags
INSERT INTO question_tag_associations (question_id, tag_id, created_at)
VALUES 
    ('aaaaaaaa-0001-aaaa-aaaa-aaaaaaaaaaaa', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW()),
    ('aaaaaaaa-0002-aaaa-aaaa-aaaaaaaaaaaa', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW()),
    ('bbbbbbbb-0001-bbbb-bbbb-bbbbbbbbbbbb', 'dddddddd-dddd-dddd-dddd-dddddddddddd', NOW()),
    ('bbbbbbbb-0002-bbbb-bbbb-bbbbbbbbbbbb', 'dddddddd-dddd-dddd-dddd-dddddddddddd', NOW()),
    ('cccccccc-0001-cccc-cccc-cccccccccccc', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NOW())
ON CONFLICT (question_id, tag_id) DO NOTHING;

-- Insert sample assessment tags
INSERT INTO assessment_tags (id, name, description, created_at)
VALUES 
    ('aaaaaaaa-1111-aaaa-aaaa-aaaaaaaaaaaa', 'go', 'Go programming assessments', NOW()),
    ('bbbbbbbb-1111-bbbb-bbbb-bbbbbbbbbbbb', 'backend', 'Backend development assessments', NOW()),
    ('cccccccc-1111-cccc-cccc-cccccccccccc', 'screening', 'Initial screening assessments', NOW()),
    ('dddddddd-1111-dddd-dddd-dddddddddddd', 'technical', 'Technical deep dive assessments', NOW())
ON CONFLICT (id) DO NOTHING;

COMMIT;
