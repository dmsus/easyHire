# API Contracts

## Authentication

### POST /api/v1/auth/login
**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}

Response (200):
json

{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJSUzI1NiIs...",
  "expires_in": 3600,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "role": "hr",
    "name": "John Doe"
  }
}
```
## Assessment Management

### POST /api/v1/assessments

**Headers:** Authorization: Bearer <token>
**Request:**
```json
{
  "title": "Go Backend Developer Assessment",
  "description": "Technical assessment for senior Go developers",
  "competencies": ["go_fundamentals", "concurrency", "system_design"],
  "target_level": "senior",
  "duration_minutes": 60,
  "question_count": 20,
  "candidates": [
    {"email": "candidate1@company.com", "name": "Alice"},
    {"email": "candidate2@company.com", "name": "Bob"}
  ],
  "deadline": "2024-12-31T23:59:59Z"
}
```
**Response (201):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "invitation_links": [
    {
      "candidate_email": "candidate1@company.com",
      "assessment_url": "https://easyhire.com/assessment/abc123"
    }
  ],
  "created_at": "2024-11-21T10:00:00Z"
}
```
## Question Generation
### POST /api/v1/questions/generate

**Headers:** Authorization: Bearer <token>

**Request:**
```json
{
  "competency": "concurrency",
  "level": "senior",
  "question_type": "coding",
  "count": 3,
  "variations": 2
}
```
**Response (200):**
```json
{
  "questions": [
    {
      "id": "q123",
      "type": "coding",
      "competency": "concurrency",
      "level": "senior",
      "content": {
        "description": "Implement a thread-safe rate limiter using goroutines and channels",
        "code_template": "package main\n\ntype RateLimiter struct {\n    // TODO: Implement\n}\n\nfunc NewRateLimiter(rate int, per time.Duration) *RateLimiter {\n    // TODO: Implement\n}",
        "test_cases": [
          {"input": "rate=10, per=1s", "expected_output": "Allows 10 requests per second"},
          {"input": "rate=100, per=1m", "expected_output": "Allows 100 requests per minute"}
        ],
        "hints": ["Use ticker for timing", "Consider using buffered channels"]
      }
    }
  ]
}
```
## Code Execution
### POST /api/v1/execute

**Headers:** Authorization: Bearer <token>

**Request:**
```json
{
  "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}",
  "language": "go",
  "timeout_ms": 5000,
  "memory_mb": 256,
  "test_cases": [
    {"input": "", "expected_output": "Hello, World!\\n"}
  ]
}
```
**Response (200):**
```json
{
  "success": true,
  "output": "Hello, World!\\n",
  "execution_time_ms": 120,
  "memory_used_mb": 45,
  "test_results": [
    {
      "passed": true,
      "input": "",
      "expected_output": "Hello, World!\\n",
      "actual_output": "Hello, World!\\n",
      "execution_time_ms": 120
    }
  ],
  "errors": []
}
```
## Results & Analytics
### GET /api/v1/assessments/{id}/results

**Headers:** Authorization: Bearer <token>

**Response (200):**
```json
{
  "assessment_id": "123e4567-e89b-12d3-a456-426614174000",
  "candidate_results": [
    {
      "candidate_id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Alice",
      "email": "candidate1@company.com",
      "overall_score": 78.5,
      "level_assessment": "senior",
      "competency_scores": {
        "go_fundamentals": 85.0,
        "concurrency": 92.0,
        "system_design": 65.0
      },
      "time_spent_minutes": 45,
      "completed_at": "2024-11-21T11:45:00Z"
    }
  ],
  "aggregate_stats": {
    "average_score": 78.5,
    "completion_rate": 100.0,
    "time_spent_avg_minutes": 45
  }
}
```
## Webhook Events
### Assessment Completed Webhook

**URL:** {client_webhook_url}/assessment/completed

**Payload:**
```json
{
  "event": "assessment.completed",
  "timestamp": "2024-11-21T11:45:00Z",
  "data": {
    "assessment_id": "123e4567-e89b-12d3-a456-426614174000",
    "candidate_id": "550e8400-e29b-41d4-a716-446655440001",
    "score": 78.5,
    "level": "senior",
    "report_url": "https://easyhire.com/reports/abc123"
  }
}
```
## Question Generated Webhook

**URL:** {client_webhook_url}/question/generated

**Payload:**
```json
{
  "event": "question.generated",
  "timestamp": "2024-11-21T10:30:00Z",
  "data": {
    "question_id": "q123",
    "competency": "concurrency",
    "level": "senior",
    "status": "pending_review",
    "expert_review_url": "https://easyhire.com/review/q123"
  }
}
```

