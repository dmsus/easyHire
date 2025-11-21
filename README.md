# EasyHire - AI-Powered Technical Assessment Platform

## Overview
EasyHire is a comprehensive B2B platform that revolutionizes technical hiring through AI-generated assessments. The platform specializes in evaluating Go developers using a sophisticated competency matrix and advanced assessment engine.

## Features
- 🤖 AI-Powered Question Generation using RAG
- 🎯 Comprehensive Go Developer Competency Assessment  
- 💻 Real-time Code Execution with Docker
- 📊 Advanced Analytics and Skill Gap Analysis
- 👥 Multi-role Platform (HR, Candidates, Technical Experts)

## Tech Stack
- **Frontend**: React, TypeScript, Tailwind CSS, Vite
- **Backend**: Go, Gin, PostgreSQL
- **AI**: Gemini API, OpenRouter, RAG
- **Infrastructure**: Docker, Kubernetes, Prometheus
- **Security**: JWT, RBAC, Secure Code Execution

## Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 14+

### Development
```bash
# Clone repository
git clone https://github.com/your-username/easyhire.git
cd easyhire

# Backend
cd backend
go mod download
go run cmd/backend/main.go

# Frontend  
cd frontend
npm install
npm run dev

# AI Service
cd ai
pip install -r requirements.txt
python main.py
```
### Project Structure

``` text
easyhire/
├── frontend/          # React frontend application
├── backend/           # Go backend API
├── ai/               # AI question generation service
├── docs/             # Documentation
├── docker/           # Docker configurations
└── k8s/              # Kubernetes manifests
```
### Contributing

Please read our Contributing Guide before submitting pull requests.

### License 

MIT License - see LICENSE file for details.