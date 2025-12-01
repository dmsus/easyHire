# 🚀 EasyHire - AI-Powered Technical Assessment Platform for Go Developers

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org)
[![React](https://img.shields.io/badge/React-18-61DAFB.svg)](https://reactjs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6.svg)](https://www.typescriptlang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![OpenAPI 3.0](https://img.shields.io/badge/OpenAPI-3.0-6BA539.svg)](https://swagger.io/specification)

**Enterprise-ready platform for automated technical assessment of Go developers using AI and competency-based evaluation.**

## ✨ Key Features

### 🤖 AI-Powered Assessment Engine
- **RAG-based Question Generation** - Contextual questions using Gemini API + OpenRouter
- **Smart Competency Matching** - 40+ Go competencies with 4-level progression (Junior → Expert)
- **Adaptive Testing** - Questions adjust to candidate's performance level
- **Automated Validation** - Technical expert review workflow for AI-generated questions

### 🎯 Advanced Evaluation System
- **Fibonacci Scoring** - Mathematical progression (1,2,3,5) for level weights
- **Real-time Code Execution** - Secure Docker containers with resource limits
- **Competency Breakdown** - Detailed skill analysis with strengths/weaknesses
- **Time-based Bonuses** - Rewards for efficient solutions (1.0-1.2x multipliers)

### 👥 Multi-Role Platform
- **HR Specialists** - Create assessments, bulk candidate invites, progress tracking
- **Candidates** - Clean test interface with Monaco editor, instant feedback
- **Technical Experts** - AI question validation, competency matrix management
- **Administrators** - System configuration, user management, analytics

### 🔒 Enterprise-Ready Security
- **JWT Authentication** with RS256 asymmetric encryption
- **Role-Based Access Control** (5 roles with granular permissions)
- **Secure Code Execution** - Docker sandboxing with time/memory limits
- **OAuth2 SSO Integration** - Auth0, Okta, Google, Microsoft Azure AD

## 🏗️ System Architecture (C4 Model)

### 📊 High-Level Overview
```text
┌─────────────────────────────────────────────────────────────┐
│ EasyHire Platform │
├─────────────┬─────────────┬────────────────┬───────────────┤
│ React │ Go │ AI │ PostgreSQL │
│ SPA │ Backend │ Service │ + pgvector │
├─────────────┼─────────────┼────────────────┼───────────────┤
│ Redis │ Code │ External │ Email/ATS │
│ Cache │ Executor │ AI Providers │ Integration │
└─────────────┴─────────────┴────────────────┴───────────────┘
```


### 🔗 External Integrations
- **AI Providers**: Gemini API, OpenRouter (Llama, Mistral)
- **HR Systems**: Greenhouse, Lever, Workable
- **SSO Providers**: Auth0, Okta, Google, Microsoft
- **Email Services**: SendGrid, AWS SES, Postmark

## 📚 Documentation

### 📋 Project Status
- **Phase 1 (Foundation)**: ✅ COMPLETED (5/5 tasks)
- **Phase 2 (Implementation)**: 🚀 READY TO START
- **Total Documentation**: 40 files created

### 📁 Key Documentation Files
```text
docs/
├── 📖 analysis/ # Requirements analysis
│ ├── user-stories.md # 16 user stories for all roles
│ └── use-cases.md # Detailed use cases with scenarios
├── 🏗️ architecture/ # System architecture
│ ├── system-architecture.md # Complete C4 model documentation
│ └── diagrams/ # 5 PlantUML C4 diagrams
├── 🎯 competency-matrix.md # 40+ Go competencies with levels
├── 📊 assessment-framework.md # Fibonacci-based scoring system
├── 🔌 api/ # Complete API specification
│ ├── openapi.yaml # OpenAPI 3.0 spec (42 endpoints)
│ ├── rest-api.md # REST API overview
│ ├── schemas/ # 5 data schemas
│ ├── paths/ # 14 endpoint definitions
│ └── examples/ # Request/response examples
└── 📈 project-status.md # Current project status
```


## 🛠️ Technology Stack

### Frontend
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS + Headless UI
- **Code Editor**: Monaco Editor (VS Code in browser)
- **State Management**: Zustand
- **Routing**: React Router v6
- **HTTP Client**: Axios + React Query

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin HTTP Framework
- **Database**: PostgreSQL 14+ with pgvector extension
- **ORM**: GORM with migrations
- **Cache**: Redis for sessions and rate limiting
- **Authentication**: JWT (RS256) + OAuth2
- **Validation**: Go Validator v10

### AI Service
- **Language**: Python 3.11+
- **Framework**: FastAPI + Pydantic
- **AI Models**: Gemini API + OpenRouter (Llama, Mistral)
- **Vector Database**: PostgreSQL + pgvector
- **Embeddings**: Sentence Transformers
- **Prompt Engineering**: LangChain + custom templates

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Orchestration**: Kubernetes (k8s manifests included)
- **Reverse Proxy**: Nginx + Let's Encrypt
- **Monitoring**: Prometheus + Grafana + Loki
- **CI/CD**: GitHub Actions with Go/Node.js/Python
- **Database Backup**: WAL-G + S3 compatible storage

## 🚀 Quick Start

### Prerequisites
```bash
# Required software
- Go 1.21+ (https://golang.org/dl/)
- Node.js 18+ (https://nodejs.org/)
- Python 3.11+ (https://www.python.org/)
- Docker & Docker Compose (https://docs.docker.com/)
- PostgreSQL 14+ (or use Docker)
```
## Development Environment Setup
### Option 1: Using Docker Compose (Recommended)
```bash

# Clone the repository
git clone https://github.com/dmsus/easyhire.git
cd easyhire

# Start all services
docker-compose -f docker/docker-compose.dev.yml up -d

# Services will be available at:
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8080
# - AI Service: http://localhost:8000
# - PostgreSQL: localhost:5432
# - Redis: localhost:6379
```
### Option 2: Manual Setup
```bash

# 1. Clone repository
git clone https://github.com/dmsus/easyhire.git
cd easyhire

# 2. Setup Backend
cd backend
cp config/.env.example config/.env
go mod download
go run cmd/backend/main.go

# 3. Setup Frontend
cd frontend
cp .env.example .env
npm install
npm run dev

# 4. Setup AI Service
cd ai
cp .env.example .env
pip install -r requirements.txt
uvicorn main:app --reload --port 8000

# 5. Setup Database
docker run -d \
  --name easyhire-postgres \
  -e POSTGRES_PASSWORD=securepassword \
  -e POSTGRES_DB=easyhire \
  -p 5432:5432 \
  postgres:14-alpine
  
# 6. Setup Redis
docker run -d \
  --name easyhire-redis \
  -p 6379:6379 \
  redis:7-alpine
```
## Configuration

Create .env files based on examples:
```bash

# Backend .env
DATABASE_URL=postgres://user:password@localhost:5432/easyhire
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-256-bit-secret
AI_SERVICE_URL=http://localhost:8000
CODE_EXECUTOR_URL=http://localhost:8081

# Frontend .env
VITE_API_URL=http://localhost:8080/v1
VITE_WS_URL=ws://localhost:8080/ws

# AI Service .env
GEMINI_API_KEY=your-gemini-api-key
OPENROUTER_API_KEY=your-openrouter-key
DATABASE_URL=postgres://user:password@localhost:5432/easyhire
```
## 📁 Project Structure
```text

easyhire/
├── 📁 backend/                 # Go backend API
│   ├── cmd/backend/           # Main application entry point
│   ├── internal/              # Private application code
│   │   ├── handlers/          # HTTP handlers (API endpoints)
│   │   ├── services/          # Business logic layer
│   │   ├── models/            # Data models (GORM structs)
│   │   ├── repository/        # Data access layer
│   │   ├── middleware/        # Authentication, logging, etc.
│   │   └── pkg/utils/         # Shared utilities
│   ├── config/                # Configuration files
│   ├── migrations/            # Database migrations
│   └── tests/                 # Integration and unit tests
├── 📁 frontend/               # React frontend
│   ├── src/
│   │   ├── components/        # Reusable React components
│   │   ├── pages/            # Page components (HR, Candidate, etc.)
│   │   ├── hooks/            # Custom React hooks
│   │   ├── services/         # API service clients
│   │   ├── stores/           # Zustand state stores
│   │   ├── utils/            # Frontend utilities
│   │   └── types/            # TypeScript type definitions
│   ├── public/               # Static assets
│   └── tests/                # Frontend tests
├── 📁 ai/                     # AI question generation service
│   ├── services/             # AI integration logic
│   ├── models/               # AI model configurations
│   ├── prompts/              # Prompt templates
│   ├── rag/                  # Retrieval-Augmented Generation
│   └── tests/                # AI service tests
├── 📁 docs/                   # Complete documentation (40 files)
│   ├── analysis/             # Requirements analysis ✅
│   ├── architecture/         # System architecture ✅
│   ├── api/                  # API specification ✅
│   ├── development/          # Development guides
│   ├── deployment/           # Deployment guides
│   └── project-status.md     # Current project status
├── 📁 docker/                # Docker configurations
│   ├── docker-compose.dev.yml  # Development environment
│   ├── docker-compose.prod.yml # Production environment
│   └── Dockerfiles/          # Individual service Dockerfiles
├── 📁 k8s/                   # Kubernetes manifests
│   ├── deployments/          # Deployment configurations
│   ├── services/             # Service definitions
│   ├── ingress/              # Ingress configurations
│   └── configmaps/           # Configuration maps
├── 📁 scripts/               # Utility scripts
├── 📄 README.md              # This file
├── 📄 LICENSE                # MIT License
├── 📄 .gitignore            # Git ignore rules
├── 📄 go.mod                # Go module definition
├── 📄 package.json          # Node.js dependencies
└── 📄 requirements.txt      # Python dependencies
```
## 📊 Assessment System
### Fibonacci-Based Scoring
```text

Level Weights (Fibonacci Sequence):
- Junior: 1 (F1)
- Middle: 2 (F2)  
- Senior: 3 (F3)
- Expert: 5 (F4)

Progression Thresholds:
- Junior: 8+ points (F6)
- Middle: 21+ points (F8)
- Senior: 55+ points (F10)
- Expert: 144+ points (F12)
```
### Competency Matrix (40+ Skills)
```yaml

Core Go Development:
  - Go Syntax & Fundamentals (weight: 1.0)
  - Data Structures (weight: 1.1)
  - Memory Management (weight: 1.1)
  - Concurrency (weight: 1.3) ⭐
  - HTTP & Web Services (weight: 1.0)

System Design:
  - System Architecture (weight: 1.3) ⭐
  - Microservices (weight: 1.2)
  - Containerization (weight: 1.1)
  - Reliability & Scalability (weight: 1.2)

Software Engineering:
  - Design Patterns (weight: 1.3) ⭐
  - Code Quality (weight: 1.2)
  - Testing & CI/CD (weight: 1.1)

Security:
  - Web Security (weight: 1.2)
  - Data Security (weight: 1.3) ⭐
```
## 🔌 API Overview
### Quick API Examples
```bash

# Create assessment (HR)
curl -X POST https://api.easyhire.com/v1/assessments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Go Developer",
    "role": "backend_developer",
    "target_level": "senior",
    "competency_weights": {
      "concurrency": 1.4,
      "system_design": 1.3
    }
  }'

# Execute Go code (Candidate)
curl -X POST https://api.easyhire.com/v1/execute \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "language": "go",
    "code": "package main\n\nfunc main() {\n    println(\"Hello, World!\")\n}",
    "test_cases": [{"input": "", "expected_output": "Hello, World!"}]
  }'

# Generate AI questions
curl -X POST https://api.easyhire.com/v1/questions/generate \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "competencies": [
      {"name": "concurrency", "level": "senior", "count": 3}
    ],
    "model": "gemini"
  }'
```
## API Features
- 42 REST Endpoints - Full CRUD for all resources
- OpenAPI 3.0 Specification - Machine-readable API docs
- JWT Authentication - With refresh token rotation
- Rate Limiting - Tier-based (60-1000 req/min)
- WebSocket Support - Real-time candidate progress
- Webhooks - 12 events with HMAC signatures
- File Downloads - PDF reports, CSV exports

## 🧪 Testing
### Test Credentials (Sandbox)
```yaml

HR Specialist:
  email: hr-test@easyhire.com
  password: Test123!

Candidate:
  email: candidate-test@easyhire.com  
  password: Test123!

Technical Expert:
  email: expert-test@easyhire.com
  password: Test123!

Sandbox API: https://sandbox.api.easyhire.com/v1
```
## Running Tests
```bash

# Backend tests
cd backend
go test ./... -v

# Frontend tests  
cd frontend
npm test

# AI service tests
cd ai
pytest

# E2E tests
cd tests/e2e
npm run test:e2e
```
## 📈 Development Roadmap
### ✅ Phase 1: Foundation (COMPLETED)
- Project structure and documentation
- User stories and use cases (16 stories)
- Competency matrix (40+ Go skills)
- Assessment framework (Fibonacci scoring)
- System architecture (C4 diagrams)
- API specification (42 endpoints, OpenAPI 3.0)

### 🚀 Phase 2: Implementation (NEXT)
- Backend API implementation (Go + Gin)
- Database design and migrations
- AI service implementation (Python + FastAPI)
- Frontend application (React + TypeScript)
- Code execution service (Docker-based)
- Authentication and authorization
- Deployment configuration

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.


## 🙏 Acknowledgments
- Google Gemini API for AI question generation
- OpenRouter for open-source model access
- C4 Model for system architecture visualization
- PlantUML for diagram generation
- Fibonacci Sequence for scoring system inspiration
