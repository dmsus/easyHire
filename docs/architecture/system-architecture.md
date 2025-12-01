# System Architecture (C4 Model)

## Level 1: System Context Diagram

![C4 Level 1](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/c4-level1.puml)

**Description:**
- **HR Specialists** create and manage technical assessments
- **Candidates** take Go programming assessments via web interface
- **Technical Experts** review and validate AI-generated questions
- **Admins** configure system settings and monitor performance
- **EasyHire Platform** coordinates assessment lifecycle with AI integration
- **External Systems** provide AI services, email notifications, and ATS integration

## Level 2: Container Diagram

![C4 Level 2](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/c4-level2.puml)

**Description:**
- **React SPA**: TypeScript frontend with role-based interfaces
- **Go Backend API**: RESTful service with Fibonacci scoring engine
- **AI Question Service**: Python service for RAG-based question generation
- **PostgreSQL + pgvector**: Primary database with vector embeddings
- **Code Execution Service**: Docker-based secure Go code runner
- **Redis**: Session management and caching layer
- **External Dependencies**: AI providers, email services, SSO providers

## Level 3: Component Diagram (Backend API)

![C4 Level 3 Backend](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/c4-level3-backend.puml)

**Description:**
- **Authentication Handler**: JWT-based authentication with RBAC
- **Assessment Handler**: Manages assessment lifecycle and candidate flow
- **Question Handler**: Question bank management and AI integration
- **Scoring Engine**: Fibonacci-based competency scoring (weights: 1,2,3,5)
- **Code Execution Client**: Secure Docker communication for code evaluation
- **Repository Layer**: GORM-based data access objects for PostgreSQL
- **External Components**: AI Service for question generation, Code Executor for safe execution

## Technology Stack

### Frontend
- **React 18 + TypeScript**: Type-safe UI development
- **Vite**: Fast build tool with hot module replacement
- **Tailwind CSS**: Utility-first CSS framework
- **Monaco Editor**: VS Code-based code editor for assessments
- **Zustand**: Lightweight state management
- **React Router v6**: Client-side routing

### Backend
- **Go 1.21**: High-performance, concurrent backend
- **Gin Framework**: Minimalist HTTP web framework
- **GORM**: ORM with struct-first approach
- **PostgreSQL 14**: Relational database with JSON support
- **Redis 7**: Caching and session storage
- **JWT (RS256)**: Secure token-based authentication

### AI Service
- **Python 3.11 + FastAPI**: AI service framework with async support
- **Gemini API**: Google's AI model for question generation
- **OpenRouter**: Fallback for open-source models (Llama, Mistral)
- **pgvector**: PostgreSQL extension for vector similarity search

### Infrastructure
- **Docker**: Containerization for services
- **Kubernetes**: Production orchestration and auto-scaling
- **Nginx**: Reverse proxy and load balancing
- **Prometheus + Grafana**: Monitoring and observability stack

## Communication Patterns
- Frontend ↔ Backend: REST API over HTTPS (JSON)
- Backend ↔ AI Service: gRPC/HTTP for question generation
- Backend ↔ Database: SQL with connection pooling
- Backend ↔ Code Executor: HTTP with timeout and resource limits
- Backend ↔ Redis: Redis protocol for caching and sessions
