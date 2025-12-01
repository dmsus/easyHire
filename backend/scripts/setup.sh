#!/bin/bash

# EasyHire Backend Setup Script
set -e

echo "🚀 Setting up EasyHire Backend..."

# Check prerequisites
echo "🔍 Checking prerequisites..."
command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed. Please install Go 1.21+"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "⚠️ Docker is not installed. Some features may not work."; }
command -v docker-compose >/dev/null 2>&1 || { echo "⚠️ Docker Compose is not installed. Some features may not work."; }

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [[ $(echo "$GO_VERSION < 1.21" | bc) -eq 1 ]]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go 1.21+"
    exit 1
fi
echo "✅ Go $GO_VERSION detected"

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download

# Setup environment
echo "⚙️ Setting up environment..."
if [ ! -f "config/.env" ]; then
    cp config/.env.example config/.env
    echo "📝 Created config/.env from example"
    echo "⚠️ Please edit config/.env with your settings"
else
    echo "✅ config/.env already exists"
fi

# Create directory for binaries
echo "📁 Creating directories..."
mkdir -p bin logs

# Setup git hooks if in git repository
if [ -d ".git" ]; then
    echo "🔗 Setting up git hooks..."
    cp -n scripts/pre-commit .git/hooks/ || true
    chmod +x .git/hooks/pre-commit 2>/dev/null || true
fi

# Build the application
echo "🔨 Building application..."
if make build; then
    echo "✅ Build successful"
else
    echo "⚠️ Build failed, but setup can continue"
fi

echo ""
echo "🎉 Setup completed!"
echo ""
echo "Next steps:"
echo "1. Edit config/.env with your settings"
echo "2. Start dependencies: make docker-up"
echo "3. Run the application: make run"
echo "4. Open http://localhost:8080 to verify"
echo ""
echo "For development with hot reload:"
echo "   go install github.com/cosmtrek/air@latest"
echo "   air"
echo ""
