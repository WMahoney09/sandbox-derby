#!/bin/bash
set -euo pipefail

echo "============================================"
echo "  Sandbox Derby — Setup"
echo "============================================"
echo ""

# Check prerequisites
echo "Checking prerequisites..."

if ! command -v docker &> /dev/null; then
    echo "  ERROR: docker is not installed"
    exit 1
fi
echo "  docker: $(docker --version)"

if ! command -v go &> /dev/null; then
    echo "  ERROR: go is not installed"
    exit 1
fi
echo "  go: $(go version)"

echo ""

# Build and install the CLI
echo "Installing derby CLI..."
go install ./cmd/derby
echo "  Installed to $(go env GOPATH)/bin/derby"

echo ""

# Set up .env if needed
if [ ! -f .env ]; then
    cp .env.example .env
    echo "Created .env from .env.example"
    echo "  Edit .env and add your ANTHROPIC_API_KEY and GITHUB_TOKEN"
else
    echo ".env already exists, skipping"
fi

echo ""

# Build the Docker image
echo "Building sandbox image..."
docker compose build
echo "  Image built: sandbox-derby:latest"

echo ""
echo "============================================"
echo "  Setup complete!"
echo "============================================"
echo ""
echo "  Next steps:"
echo "    1. Edit .env with your API keys (if you haven't already)"
echo "    2. derby drive                     — interactive sandbox"
echo "    3. derby coast --help              — autonomous sandbox"
echo "    4. derby scrimmage examples/*.yaml  — run a scrimmage"
echo ""
