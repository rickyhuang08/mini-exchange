#!/bin/bash

set -e

echo "🚀 Installing Go dependencies..."

# Check if go is installed
if ! command -v go &> /dev/null
then
    echo "❌ Go is not installed!"
    exit 1
fi

# Initialize go module if not exists
if [ ! -f "go.mod" ]; then
    echo "📦 go.mod not found. Initializing module..."
    read -p "Enter module name (e.g. github.com/yourname/project): " module_name
    go mod init "$module_name"
fi

echo "📥 Downloading dependencies..."

go get github.com/google/uuid
go get github.com/gorilla/websocket
go get golang.org/x/crypto/bcrypt
go get github.com/joho/godotenv
go get gopkg.in/yaml.v3
go get github.com/gin-gonic/gin@v1.11.0
go get github.com/golang-jwt/jwt/v5

echo "🧹 Running go mod tidy..."
go mod tidy

echo "✅ All dependencies installed successfully!"