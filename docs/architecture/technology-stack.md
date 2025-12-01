# Technology Stack Decisions

## Frontend
- **React 18 + TypeScript**: Type-safe UI development with modern features
- **Vite**: Fast builds, hot reload, optimized production bundles
- **Tailwind CSS**: Utility-first CSS for rapid, consistent UI development
- **React Router v6**: Declarative routing with nested routes support
- **Zustand**: Lightweight state management without boilerplate
- **Monaco Editor**: VS Code-based code editor for assessments
- **Axios**: Promise-based HTTP client for API calls

## Backend
- **Go 1.21**: High performance, excellent concurrency, static typing
- **Gin Web Framework**: Minimalistic, fast HTTP web framework
- **GORM**: ORM with struct-first approach and migrations
- **PostgreSQL 14**: Reliable, feature-rich relational database
- **Redis 7**: In-memory data store for caching and sessions
- **JWT (RS256)**: Secure token-based authentication
- **Bcrypt**: Password hashing with adjustable cost
- **Docker SDK**: Programmatic container management for code execution

## AI Service
- **Python 3.11**: Rich AI/ML ecosystem and rapid development
- **FastAPI**: Modern, fast (Starlette) web framework with automatic docs
- **Gemini API**: Google's advanced AI model (free tier available)
- **OpenRouter**: Access to multiple open-source models (Llama, Mistral)
- **pgvector**: PostgreSQL extension for vector similarity search
- **Sentence Transformers**: For generating text embeddings
- **Pydantic**: Data validation and settings management

## Code Execution
- **Docker**: Containerization for secure, isolated code execution
- **Resource constraints**: CPU, memory, process, and network limits
- **Read-only filesystems**: Except designated temporary directories
- **User namespace**: Non-root execution inside containers
- **Timeouts**: Strict execution time limits per question

## Infrastructure & DevOps
- **Docker & Docker Compose**: Containerization and local development
- **Kubernetes**: Production orchestration and auto-scaling
- **Nginx**: Reverse proxy, load balancing, SSL termination
- **Prometheus & Grafana**: Metrics collection and visualization
- **GitHub Actions**: CI/CD pipelines with automated testing
- **PostgreSQL with TimescaleDB**: For time-series metrics if needed
- **AWS/GCP/Azure**: Cloud-agnostic design for flexibility

## Security
- **JWT with RS256**: Asymmetric encryption for tokens
- **Role-Based Access Control (RBAC)**: Fine-grained permissions
- **Rate Limiting**: Per-user and per-IP request limits
- **Input Validation**: Comprehensive validation at all layers
- **Security Headers**: CSP, HSTS, XSS protection
- **Docker Security**: Seccomp profiles, AppArmor, capabilities dropping

## Monitoring & Observability
- **Structured Logging**: JSON logs with correlation IDs
- **OpenTelemetry**: Distributed tracing across services
- **Health Checks**: Readiness and liveness endpoints
- **Custom Metrics**: Business metrics (assessments, questions, users)
- **Alerting**: Proactive monitoring with Alertmanager
