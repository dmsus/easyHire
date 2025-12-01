# System Architecture (C4 Model)

## Level 1: System Context Diagram

![C4 Level 1](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/yourusername/easyhire/main/docs/architecture/diagrams/c4-level1-context.puml)

**Description:**
EasyHire connects HR specialists, candidates, and technical experts for Go developer assessments with Fibonacci-based scoring.

**Actors:**
- **HR Specialists**: Create assessments using competency matrix
- **Candidates**: Take Go assessments with real-time scoring
- **Technical Experts**: Validate AI-generated questions
- **System Admins**: Configure platform settings

**External Systems:**
- **AI Providers**: Gemini API, OpenRouter for question generation
- **Email Services**: For notifications and invitations
- **ATS Systems**: Applicant Tracking System integration
- **SSO Providers**: Enterprise authentication

## Level 2: Container Diagram

![C4 Level 2](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/yourusername/easyhire/main/docs/architecture/diagrams/c4-level2-containers.puml)

**Containers:**
1. **React SPA**: User interfaces for all roles
2. **Go Backend API**: Business logic and Fibonacci scoring
3. **AI Question Service**: Python service for question generation
4. **PostgreSQL + pgvector**: Database with vector embeddings
5. **Code Execution Service**: Secure Docker containers for Go code
6. **Redis**: Caching and session management

**Communication:**
- Frontend ↔ Backend: REST API
- Backend ↔ AI Service: gRPC/HTTP
- Backend ↔ Database: SQL
- Backend ↔ Code Executor: HTTP
- Backend ↔ Redis: Redis protocol

## Level 3: Component Diagrams

### Backend API Components
![C4 Level 3 Backend](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/yourusername/easyhire/main/docs/architecture/diagrams/c4-level3-backend.puml)

**Key Components:**
- **Fibonacci Scoring Engine**: Calculates scores using weights (1,2,3,5)
- **Competency Evaluator**: Maps answers to 40+ Go competencies
- **Assessment Orchestrator**: Manages 60-minute assessments
- **Code Execution Client**: Communicates with Docker service

### Frontend Application Components
![C4 Level 3 Frontend](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/yourusername/easyhire/main/docs/architecture/diagrams/c4-level3-frontend.puml)

**UI Components:**
- **HR Dashboard**: Assessment creation and candidate management
- **Candidate Interface**: Code editor with real-time feedback
- **Expert Panel**: Question validation and review
- **Admin Console**: System configuration

### AI Service Components
![C4 Level 3 AI](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/yourusername/easyhire/main/docs/architecture/diagrams/c4-level3-ai.puml)

**AI Components:**
- **RAG Engine**: Retrieval-Augmented Generation with vector search
- **Question Generator**: Creates questions based on competencies
- **Quality Validator**: Ensures question quality and correctness

## Technology Stack Summary

**Frontend:** React 18, TypeScript, Vite, Tailwind CSS  
**Backend:** Go 1.21, Gin, GORM, PostgreSQL  
**AI Service:** Python 3.11, FastAPI, Gemini API, OpenRouter  
**Infrastructure:** Docker, Kubernetes, Redis, pgvector  
**Security:** JWT, RBAC, Docker isolation, TLS 1.3

## Data Flow Overview

1. HR creates assessment → AI generates questions → Candidates take test  
2. Code execution → Fibonacci scoring → Results analysis  
3. Expert validation → AI learning → Question improvement  

*(Detailed API specifications will be created in Task #5)*
