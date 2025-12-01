# Security Architecture Design

## Authentication & Authorization

### JWT-Based Authentication

┌─────────────┐ 1. Login ┌─────────────┐
│ Client │ ──────────────► │ Backend │
│ │ │ │
│ │ ◄────────────── │ │
│ │ 2. JWT Token │ │
└─────────────┘ └─────────────┘

**Implementation:**
- **Algorithm**: RS256 (asymmetric) for access tokens
- **Claims**: user_id, role, permissions, exp, iat
- **Refresh tokens**: Short-lived (1h) access, long-lived (7d) refresh
- **Token storage**: HTTP-only cookies for web, secure storage for mobile

### Role-Based Access Control (RBAC)
```yaml
Roles:
  - candidate:     [take_assessment, view_own_results]
  - hr:            [create_assessment, view_candidates, export_results]
  - technical_expert: [review_questions, manage_competencies]
  - admin:         [manage_users, view_logs, configure_system]

Permissions are hierarchical and composable
```
## Data Protection

### Encryption Strategy

┌─────────────────┐    AES-256-GCM    ┌─────────────────┐
│   Sensitive     │ ────────────────► │   Database      │
│     Data        │   encrypted at    │                 │
│   (e.g., PII)   │      rest         │                 │
└─────────────────┘                   └─────────────────┘
        │                                      │
        │ TLS 1.3                              │ TLS 1.3
        ▼                                      ▼
┌─────────────────┐                   ┌─────────────────┐
│   Application   │                   │   Backup        │
│     Layer       │                   │    Storage      │
└─────────────────┘                   └─────────────────┘

### Protected Data:
- Passwords: bcrypt with work factor 12
- Personal data: AES-256-GCM with key rotation
- Assessment answers: Encrypted in transit and at rest
- API keys: Hashed with salt

## Code Execution Security

### Docker Sandbox Architecture
┌─────────────────────────────────────────────────┐
│            Host Machine                         │
│  ┌──────────────┐ ┌──────────────┐            │
│  │  Container   │ │  Container   │            │
│  │    (Go)      │ │  (Python)    │            │
│  │              │ │              │            │
│  │  CPU Limit   │ │  CPU Limit   │            │
│  │  Memory Limit│ │  Memory Limit│            │
│  │  No Network  │ │  No Network  │            │
│  │  Read-only   │ │  Read-only   │            │
│  │  No Root     │ │  No Root     │            │
│  └──────────────┘ └──────────────┘            │
│          │               │                    │
│          └───────┬───────┘                    │
│                  ▼                            │
│         ┌────────────────┐                    │
│         │  Docker Daemon │                    │
│         │  with Seccomp  │                    │
│         │   AppArmor     │                    │
│         └────────────────┘                    │
└─────────────────────────────────────────────────┘
### Container Restrictions:

- **User:** Non-root user inside container (uid 1000)
- **Capabilities:** Drop all, add only NET_BIND_SERVICE if needed
- **Seccomp:** Strict profile blocking dangerous syscalls
- **AppArmor:** Deny write except /tmp, deny network, deny process creation
- **Resources:** CPU quota, memory limit, process limit
- **Filesystem:** Read-only root, tmpfs for /tmp

## API Security

### Rate Limiting Strategy
```go 
type RateLimiter struct {
    // Token bucket algorithm
    UserLimits: map[string]RateLimit{
        "candidate": {Requests: 100, Window: "1h"},
        "hr":        {Requests: 1000, Window: "1h"},
        "api_key":   {Requests: 10000, Window: "1h"},
    }
    
    // IP-based limits for anonymous endpoints
    IPLimits: {Requests: 50, Window: "1m"}
}
```
### Protections:
- **Input Validation:** All inputs validated with strict schemas
- **SQL Injection:** Parameterized queries only (GORM)
- **XSS Prevention:** Content-Type headers, output encoding
- **CSRF Protection:** SameSite cookies, anti-CSRF tokens
- **CORS:** Strict origin validation, preflight caching

## Infrastructure Security

### Network Segmentation

┌─────────────────────────────────────────────────┐
│              Public Internet                    │
│                    │                            │
│                    ▼                            │
│          ┌──────────────────┐                  │
│          │   Load Balancer  │                  │
│          │    (nginx)       │                  │
│          │  - SSL/TLS       │                  │
│          │  - WAF           │                  │
│          │  - Rate Limiting │                  │
│          └──────────────────┘                  │
│                    │                            │
│                    ▼                            │
│    ┌────────────────────────────────────┐     │
│    │         DMZ Network                │     │
│    │  ┌─────────┐ ┌─────────┐          │     │
│    │  │  Front- │ │  API    │          │     │
│    │  │   end   │ │ Gateway │          │     │
│    │  └─────────┘ └─────────┘          │     │
│    └────────────────────────────────────┘     │
│                    │                            │
│                    ▼                            │
│    ┌────────────────────────────────────┐     │
│    │      Application Network           │     │
│    │  ┌─────────┐ ┌─────────┐          │     │
│    │  │ Backend │ │  AI     │          │     │
│    │  │   API   │ │ Service │          │     │
│    │  └─────────┘ └─────────┘          │     │
│    └────────────────────────────────────┘     │
│                    │                            │
│                    ▼                            │
│    ┌────────────────────────────────────┐     │
│    │        Database Network            │     │
│    │  ┌─────────┐ ┌─────────┐          │     │
│    │  │   Main  │ │ Vector  │          │     │
│    │  │   DB    │ │   DB    │          │     │
│    │  └─────────┘ └─────────┘          │     │
│    └────────────────────────────────────┘     │
└─────────────────────────────────────────────────┘

## Monitoring & Compliance

### Security Monitoring
- **Audit Logging:** All authentication, authorization, data access
- **SIEM Integration:** Log aggregation and anomaly detection
- **Regular Scans:** DAST/SAST weekly, penetration testing quarterly
- **Compliance:** GDPR, CCPA, SOC2 readiness

### Incident Response
1. **Detection:** Automated alerts for suspicious patterns
2. **Containment:** Automatic IP blocking, session termination
3. **Investigation:** Full audit trail and forensics
4. **Recovery:** Point-in-time recovery, data integrity checks
5. **Post-Mortem:** Root cause analysis and process improvement

