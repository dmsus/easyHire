# EasyHire REST API Documentation

## 📋 Overview

The EasyHire API follows RESTful principles and provides programmatic access to the AI-powered technical assessment platform. All API endpoints are versioned and require authentication.

## 🔐 Authentication

### JWT Bearer Tokens (Primary)

Authorization: Bearer <access_token>

**Token Types:**
- Access tokens (short-lived, 15 minutes)
- Refresh tokens (long-lived, 7 days)
- API keys (for service-to-service communication)

### OAuth 2.0 (Enterprise SSO)
Supported providers: Auth0, Okta, Google, Microsoft Azure AD

### Role-Based Access Control
| Role | Permissions |
|------|-------------|
| `candidate` | Take assessments, view own results |
| `hr` | Create assessments, manage candidates, view results |
| `expert` | Review questions, manage competency matrix |
| `admin` | System configuration, user management |
| `system` | Internal service communication |

## 📊 Rate Limiting

| Tier | Requests/Minute | Burst |
|------|----------------|-------|
| Free | 60 | 100 |
| Pro | 300 | 500 |
| Enterprise | 1000 | 2000 |

Headers included in all responses:

X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
X-RateLimit-Reset: 1637830200


## 🔄 Pagination

All list endpoints support pagination:

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20, max: 100)
- `sort` (field to sort by)
- `order` (asc/desc)

**Response Headers:**

X-Total-Count: 150
X-Page-Count: 8
X-Current-Page: 1
X-Per-Page: 20



## 📍 Endpoints

### Authentication
- `POST /auth/login` - Email/password authentication
- `POST /auth/refresh` - Refresh JWT tokens
- `POST /auth/logout` - Invalidate tokens
- `GET /auth/profile` - Get current user profile

### Assessments
- `GET /assessments` - List assessments (filtered by role)
- `POST /assessments` - Create new assessment
- `GET /assessments/{id}` - Get assessment details
- `PUT /assessments/{id}` - Update assessment
- `DELETE /assessments/{id}` - Delete assessment
- `POST /assessments/{id}/invite` - Invite candidates
- `GET /assessments/{id}/candidates` - List invited candidates
- `POST /assessments/{id}/publish` - Publish assessment
- `POST /assessments/{id}/clone` - Clone assessment

### Candidates
- `GET /candidates` - List candidates (HR only)
- `POST /candidates` - Create candidate record
- `GET /candidates/{id}` - Get candidate profile
- `PUT /candidates/{id}` - Update candidate
- `GET /candidates/{id}/assessments` - Get candidate's assessments
- `POST /candidates/{id}/assessments/{assessmentId}/start` - Start assessment
- `POST /candidates/{id}/assessments/{assessmentId}/submit` - Submit assessment
- `GET /candidates/{id}/assessments/{assessmentId}/progress` - Get progress

### Questions
- `GET /questions` - List questions (with filters)
- `POST /questions` - Create manual question
- `POST /questions/generate` - Generate AI questions
- `GET /questions/{id}` - Get question details
- `PUT /questions/{id}` - Update question
- `POST /questions/{id}/validate` - Validate AI-generated question
- `DELETE /questions/{id}` - Archive question

### Code Execution
- `POST /execute` - Execute code with test cases
- `POST /execute/test` - Test code execution (dry run)
- `GET /execute/{id}/status` - Get execution status
- `GET /execute/{id}/output` - Get execution output

### Results
- `GET /results` - List results (filtered by role)
- `GET /results/{id}` - Get detailed results
- `GET /results/{id}/breakdown` - Get competency breakdown
- `GET /results/{id}/pdf` - Download PDF report
- `POST /results/{id}/share` - Share results via email
- `GET /results/{id}/comparison` - Compare with role requirements

### Webhooks
- `GET /webhooks` - List registered webhooks
- `POST /webhooks` - Create webhook
- `PUT /webhooks/{id}` - Update webhook
- `DELETE /webhooks/{id}` - Delete webhook
- `POST /webhooks/{id}/test` - Test webhook

## 🚦 Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request succeeded |
| 201 | Created - Resource created |
| 202 | Accepted - Request accepted for processing |
| 204 | No Content - Success, no response body |
| 400 | Bad Request - Invalid parameters |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource conflict |
| 422 | Unprocessable Entity - Validation failed |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Server error |
| 502 | Bad Gateway - Upstream service error |
| 503 | Service Unavailable - Maintenance or overload |

## 🎯 Error Handling

All errors follow the same format:
```json
{
  "error": true,
  "code": "ERROR_CODE",
  "message": "Human-readable message",
  "details": { /* Additional details */ },
  "request_id": "req_123456789",
  "timestamp": "2024-01-15T10:30:00Z"
}
```
Common error codes:
- VALIDATION_ERROR - Request validation failed
- AUTH_REQUIRED - Authentication required
- INSUFFICIENT_PERMISSIONS - Role-based access denied
- RESOURCE_NOT_FOUND - Requested resource doesn't exist
- RATE_LIMIT_EXCEEDED - Too many requests
- ASSESSMENT_COMPLETED - Assessment already completed
- ASSESSMENT_EXPIRED - Assessment time expired
- EXECUTION_TIMEOUT - Code execution timeout
- AI_SERVICE_UNAVAILABLE - AI generation service down

## 🔗 Webhooks
### Events
- assessment.created
- assessment.published
- candidate.invited
- candidate.started
- candidate.completed
- candidate.expired
- question.generated
- question.validated
- code.executed
- result.calculated
- result.ready

### Payload Format
```json

{
  "event": "candidate.completed",
  "data": {
    "assessment_id": "123e4567-e89b-12d3-a456-426614174000",
    "candidate_id": "123e4567-e89b-12d3-a456-426614174001",
    "score": 78.5,
    "completed_at": "2024-01-15T10:30:00Z"
  },
  "webhook_id": "wh_123456789",
  "timestamp": "2024-01-15T10:30:00Z",
  "signature": "sha256=..."
}
```
## 📡 Real-time Updates (WebSocket)

Connect to: wss://api.easyhire.com/v1/ws

### Events:
- assessment_progress - Candidate progress updates
- code_execution_status - Code execution status
- question_generation_status - AI generation progress
- system_notification - System notifications

### 📚 SDKs & Libraries
- **Go SDK:** go get github.com/easyhire/easyhire-go
- **JavaScript SDK:** npm install @easyhire/sdk
- **Python SDK:** pip install easyhire
- **Postman Collection:** Available in /docs/postman/

## 🔍 Testing

### Sandbox Environment:
```text

https://sandbox.api.easyhire.com/v1
```
### Test Credentials:
- **HR:** hr-test@easyhire.com / Test123!
- **Candidate:** candidate-test@easyhire.com / Test123!
- **Expert:** expert-test@easyhire.com / Test123!

## 📞 Support
- API Status: https://status.easyhire.com
- Documentation: https://docs.easyhire.com
- Support: api-support@easyhire.com
- Bug Reports: GitHub Issues