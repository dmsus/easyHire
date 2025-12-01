# EasyHire Backend API

Go backend for the EasyHire AI-powered technical assessment platform.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (recommended)

### Development with Docker (Recommended)
```bash
# Start dependencies
make docker-up

# Copy environment file
cp config/.env.example config/.env

# Run the application
make run

# Or for hot reload
air
```
## Manual Setup
```bash

# Install dependencies
go mod download

# Setup database
createdb easyhire
psql easyhire -f ../docker/init-db.sql

# Copy and configure environment
cp config/.env.example config/.env
# Edit config/.env with your settings

# Run migrations (when available)
make migrate

# Start server
make run
```
## 📁 Project Structure
```text

backend/
├── cmd/backend/           # Main application entry point
├── internal/              # Private application code
│   ├── handlers/         # HTTP handlers (API endpoints)
│   ├── services/         # Business logic layer
│   ├── models/           # Data models (GORM structs)
│   ├── repository/       # Data access layer
│   ├── middleware/       # HTTP middleware
│   └── pkg/              # Internal shared packages
├── pkg/                  # Public shared packages
├── api/                  # API definitions (OpenAPI, gRPC)
├── config/               # Configuration files
├── migrations/           # Database migrations
├── tests/                # Tests
└── scripts/              # Utility scripts
```
## 🛠️ Development Commands
```bash

# Build application
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make format

# Lint code
make lint

# Tidy dependencies
make tidy

# Run all checks
make check

# Start development environment
make dev
```
## 🔧 Configuration

Configuration is loaded from environment variables via config/.env file.

Create your configuration:
```bash

cp config/.env.example config/.env
# Edit config/.env with your settings
```
Required environment variables:
- DB_* - Database configuration
- REDIS_* - Redis configuration
- JWT_SECRET - JWT signing secret
- AI_SERVICE_URL - AI service endpoint
- CODE_EXECUTOR_URL - Code execution service endpoint

## 📡 API Endpoints
### Health Checks
- GET / - API information
- GET /health - Health status with dependencies
- GET /ready - Readiness probe
- GET /live - Liveness probe
- GET /metrics - Prometheus metrics

### API v1
- POST /api/v1/auth/login - User authentication
- POST /api/v1/auth/refresh - Refresh tokens
- POST /api/v1/auth/logout - User logout

(More endpoints to be implemented)
## 🐳 Docker Development

Start all services:
```bash

make docker-up
```
Stop services:
```bash

make docker-down
```
View logs:
```bash

docker-compose -f ../docker/docker-compose.dev.yml logs -f
```
## 🧪 Testing

Run all tests:
```bash

make test
```
Run specific test:
```bash

go test ./internal/handlers -v
```
Run with coverage:
```bash

make test-coverage
```
## 📊 Monitoring
- **Health:** http://localhost:8080/health
- **Metrics:** http://localhost:8080/metrics (Prometheus format)
- **Readiness:** http://localhost:8080/ready
- **Liveness:** http://localhost:8080/live

## 🔐 Security
- JWT authentication with RS256
- CORS configured via environment variables
- Rate limiting middleware
- Input validation
- SQL injection prevention via GORM
- Secure password hashing

## 📈 Performance
- Connection pooling for PostgreSQL and Redis
- Query optimization with GORM
- Redis caching for frequent queries
- Compression middleware
- Prometheus metrics collection

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run all checks (make check)
6. Submit a pull request

## 📄 License

MIT - See LICENSE file for details.