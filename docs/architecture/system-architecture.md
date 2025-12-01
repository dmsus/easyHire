# System Architecture (C4 Model)

## Level 1: System Context Diagram

![C4 Level 1](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/level1-context.puml)

**Description:**
EasyHire is an AI-powered platform for assessing Go developers. The system connects HR specialists, candidates, technical experts, and administrators with external AI services and email systems.

## Level 2: Container Diagram

![C4 Level 2](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/level2-containers.puml)

**Containers:**
1. **React SPA** - User interface for all roles
2. **Go Backend API** - Core business logic and assessment management
3. **AI Question Service** - Generates questions using AI models
4. **PostgreSQL + pgvector** - Primary database with vector embeddings
5. **Code Execution Service** - Secure Docker containers for Go code execution
6. **Redis** - Caching and session management

## Level 3: Component Diagrams

### Backend API Components
![C4 Level 3 Backend](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/level3-backend.puml)

**Key Components:**
- Authentication Handler - JWT-based authentication
- Assessment Handler - Manages assessment lifecycle
- Question Handler - Question bank management
- Scoring Engine - Fibonacci-based scoring calculations
- Code Execution Client - Secure Docker communication

### Frontend Application Components
![C4 Level 3 Frontend](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/level3-frontend.puml)

**UI Components:**
- HR Dashboard - Assessment creation and candidate management
- Candidate Test Interface - Code editor with real-time execution
- Expert Panel - Question validation and review
- Admin Console - System configuration

### AI Service Components
![C4 Level 3 AI](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/level3-ai.puml)

**AI Components:**
- Generation Handler - Handles question generation requests
- RAG Engine - Retrieval-Augmented Generation with vector search
- Prompt Engine - Creates optimized prompts for different question types

## Technology Stack

### Frontend
- **React 18 + TypeScript** - Type-safe UI development
- **Vite** - Fast build tool and development server
- **Tailwind CSS** - Utility-first CSS framework
- **Monaco Editor** - VS Code-based code editor
- **Zustand** - Lightweight state management
- **React Router v6** - Client-side routing

### Backend
- **Go 1.21** - High-performance backend with excellent concurrency
- **Gin Framework** - Minimalist HTTP web framework
- **GORM** - ORM for database interactions
- **PostgreSQL 14** - Relational database with JSON support
- **Redis** - Caching and session storage
- **JWT** - Token-based authentication

### AI Service
- **Python 3.11 + FastAPI** - AI/ML service framework
- **Gemini API** - Primary AI model for question generation
- **OpenRouter** - Fallback for open-source models
- **pgvector** - Vector similarity search extension

### Infrastructure
- **Docker** - Containerization for services
- **Kubernetes** - Production orchestration
- **Nginx** - Reverse proxy and load balancing

## Communication Patterns
- Frontend ↔ Backend: REST API over HTTPS
- Backend ↔ AI Service: gRPC/HTTP
- Backend ↔ Database: SQL with connection pooling
- Backend ↔ Code Executor: HTTP with timeout handling
