# System Architecture (C4 Model)

## Level 1: System Context Diagram

![C4 Level 1](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/c4-level1.puml)

**Description:**
EasyHire connects HR specialists, candidates, technical experts, and administrators with AI-powered assessment platform.

## Level 2: Container Diagram

![C4 Level 2](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/c4-level2.puml)

**Containers:**
- **React SPA**: User interface for all roles
- **Go Backend API**: Core business logic and assessment management
- **AI Question Service**: Generates questions using AI models
- **PostgreSQL + pgvector**: Primary database with vector embeddings
- **Code Execution Service**: Secure Docker containers for Go code execution
- **Redis**: Caching and session management

## Level 3: Component Diagrams

### Backend API Components
![C4 Level 3 Backend](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/c4-level3-backend.puml)

**Components:**
- Authentication Handler
- Assessment Handler  
- Question Handler
- Scoring Engine
- Code Execution Client
- Data repositories for users, assessments, questions, results

### Frontend Application Components
![C4 Level 3 Frontend](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/c4-level3-frontend.puml)

**Components:**
- HR Dashboard
- Candidate Test Interface
- Expert Panel
- Admin Console
- API Client, Code Editor, State Store, Router

### AI Service Components
![C4 Level 3 AI](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/dmsus/easyHire/main/docs/architecture/diagrams/c4-level3-ai.puml)

**Components:**
- Generation Handler
- Validation Handler
- RAG Engine
- Prompt Engine
- Vector database integration

## Technology Stack Decisions

### Frontend
- **React 18 + TypeScript**: Type safety for complex assessment UI
- **Vite**: Fast development and optimized production builds
- **Tailwind CSS**: Rapid UI development with consistent design
- **Monaco Editor**: Industry-standard code editing experience
- **Zustand**: Simple state management without boilerplate
- **React Router v6**: Declarative routing with nested routes

### Backend
- **Go 1.21**: Native Go execution, excellent performance for concurrent assessments
- **Gin Framework**: Minimalist HTTP framework with middleware support
- **GORM**: ORM with struct-first approach and migrations
- **PostgreSQL 14**: Relational database with JSON and vector extensions
- **Redis**: Low-latency caching for sessions and assessment data
- **JWT (RS256)**: Secure token-based authentication with asymmetric encryption

### AI Service
- **Python 3.11 + FastAPI**: Rich ML ecosystem and async capabilities
- **Gemini API**: High-quality question generation (free tier available)
- **OpenRouter**: Fallback to open-source models (Llama, Mistral)
- **pgvector**: Vector similarity search within PostgreSQL
- **Sentence Transformers**: For generating text embeddings

### Infrastructure
- **Docker**: Containerization for consistent environments
- **Kubernetes**: Orchestration for scalability and resilience
- **Nginx**: Reverse proxy, load balancing, SSL termination
- **Prometheus + Grafana**: Metrics collection and visualization

## API Contracts & Data Flow

### Primary Data Flows
1. **Assessment Creation**: HR → Frontend → Backend → AI Service → Database
2. **Candidate Testing**: Candidate → Frontend → Backend → Code Execution → Results
3. **Question Generation**: Backend → AI Service → Vector DB → External AI → Database

### Key API Endpoints

POST /api/v1/assessments # Create assessment
GET /api/v1/assessments/{id} # Get assessment
POST /api/v1/questions/generate # Generate AI questions
POST /api/v1/execute # Execute code
GET /api/v1/results/{id} # Get assessment results
text


### Communication Patterns
- **Frontend ↔ Backend**: REST over HTTPS with JWT authentication
- **Backend ↔ AI Service**: gRPC/HTTP for low-latency question generation
- **Backend ↔ Database**: Connection-pooled SQL queries
- **Backend ↔ Code Execution**: HTTP with strict timeouts and resource limits
- **Backend ↔ Redis**: Redis protocol for caching and session storage
